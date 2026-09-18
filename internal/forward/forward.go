package forward

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xJacky/Nginx-UI/internal/firewall"
	"github.com/0xJacky/Nginx-UI/model"
	"github.com/uozi-tech/cosy/logger"
)

type ForwardInstance struct {
	Rule       *model.PortForwardRule
	tcpListener net.Listener
	udpConn     *net.UDPConn
	stopChan    chan struct{}
	ActiveConns int64
	RxBytes     uint64
	TxBytes     uint64
}

type Manager struct {
	instances map[uint64]*ForwardInstance
	mu        sync.RWMutex
}

var GlobalManager = &Manager{
	instances: make(map[uint64]*ForwardInstance),
}

// StartRule starts a forwarding listener for the specified rule.
func (m *Manager) StartRule(rule *model.PortForwardRule) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Stop existing instance if already running
	if inst, exists := m.instances[rule.ID]; exists {
		close(inst.stopChan)
		if inst.tcpListener != nil {
			_ = inst.tcpListener.Close()
		}
		if inst.udpConn != nil {
			_ = inst.udpConn.Close()
		}
		delete(m.instances, rule.ID)
	}

	if !rule.Enabled {
		return nil
	}

	inst := &ForwardInstance{
		Rule:     rule,
		stopChan: make(chan struct{}),
		RxBytes:  rule.RxBytes,
		TxBytes:  rule.TxBytes,
	}

	listenAddr := fmt.Sprintf("%s:%d", rule.ListenIP, rule.ListenPort)
	targetAddr := fmt.Sprintf("%s:%d", rule.TargetIP, rule.TargetPort)

	// TCP
	if rule.Protocol == "tcp" || rule.Protocol == "both" || rule.Protocol == "tcp/udp" {
		ln, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return fmt.Errorf("failed to listen TCP on %s: %v", listenAddr, err)
		}
		inst.tcpListener = ln
		go m.handleTCP(inst, targetAddr)
	}

	// UDP
	if rule.Protocol == "udp" || rule.Protocol == "both" || rule.Protocol == "tcp/udp" {
		uAddr, err := net.ResolveUDPAddr("udp", listenAddr)
		if err != nil {
			if inst.tcpListener != nil {
				_ = inst.tcpListener.Close()
			}
			return fmt.Errorf("failed to resolve UDP on %s: %v", listenAddr, err)
		}
		uConn, err := net.ListenUDP("udp", uAddr)
		if err != nil {
			if inst.tcpListener != nil {
				_ = inst.tcpListener.Close()
			}
			return fmt.Errorf("failed to listen UDP on %s: %v", listenAddr, err)
		}
		inst.udpConn = uConn
		go m.handleUDP(inst, targetAddr)
	}

	m.instances[rule.ID] = inst

	// Auto open firewall port if requested
	if rule.AutoOpenFirewall {
		go func(r *model.PortForwardRule) {
			_, _ = firewall.OpenPort(fmt.Sprintf("%d", r.ListenPort), r.Protocol, "any", "Auto-opened for forward rule: "+r.Name)
		}(rule)
	}

	logger.Infof("Port forward rule [%s] started: %s -> %s (%s)", rule.Name, listenAddr, targetAddr, rule.Protocol)
	return nil
}

func (m *Manager) StopRule(ruleID uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if inst, exists := m.instances[ruleID]; exists {
		close(inst.stopChan)
		if inst.tcpListener != nil {
			_ = inst.tcpListener.Close()
		}
		if inst.udpConn != nil {
			_ = inst.udpConn.Close()
		}

		// Persist stats back to database
		db := model.UseDB()
		_ = db.Model(&model.PortForwardRule{}).Where("id = ?", ruleID).Updates(map[string]interface{}{
			"rx_bytes": inst.RxBytes,
			"tx_bytes": inst.TxBytes,
		})

		delete(m.instances, ruleID)
	}
}

func (m *Manager) handleTCP(inst *ForwardInstance, targetAddr string) {
	for {
		select {
		case <-inst.stopChan:
			return
		default:
		}

		clientConn, err := inst.tcpListener.Accept()
		if err != nil {
			return
		}

		atomic.AddInt64(&inst.ActiveConns, 1)

		go func(c net.Conn) {
			defer func() {
				_ = c.Close()
				atomic.AddInt64(&inst.ActiveConns, -1)
			}()

			targetConn, err := net.DialTimeout("tcp", targetAddr, 5*time.Second)
			if err != nil {
				return
			}
			defer targetConn.Close()

			var wg sync.WaitGroup
			wg.Add(2)

			// Client -> Target
			go func() {
				defer wg.Done()
				n, _ := copyWithCounter(targetConn, c, &inst.RxBytes)
				atomic.AddUint64(&inst.RxBytes, uint64(n))
			}()

			// Target -> Client
			go func() {
				defer wg.Done()
				n, _ := copyWithCounter(c, targetConn, &inst.TxBytes)
				atomic.AddUint64(&inst.TxBytes, uint64(n))
			}()

			wg.Wait()
		}(clientConn)
	}
}

func (m *Manager) handleUDP(inst *ForwardInstance, targetAddr string) {
	buf := make([]byte, 65535)
	targetUDPAddr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		logger.Errorf("Resolve target UDP failed: %v", err)
		return
	}

	for {
		select {
		case <-inst.stopChan:
			return
		default:
		}

		n, clientAddr, err := inst.udpConn.ReadFromUDP(buf)
		if err != nil {
			return
		}

		atomic.AddUint64(&inst.RxBytes, uint64(n))

		// Send to target
		go func(data []byte, cAddr *net.UDPAddr) {
			relayConn, err := net.DialUDP("udp", nil, targetUDPAddr)
			if err != nil {
				return
			}
			defer relayConn.Close()

			_ = relayConn.SetDeadline(time.Now().Add(10 * time.Second))
			_, _ = relayConn.Write(data)

			respBuf := make([]byte, 65535)
			respN, _, err := relayConn.ReadFrom(respBuf)
			if err == nil && respN > 0 {
				_, _ = inst.udpConn.WriteToUDP(respBuf[:respN], cAddr)
				atomic.AddUint64(&inst.TxBytes, uint64(respN))
			}
		}(append([]byte(nil), buf[:n]...), clientAddr)
	}
}

func copyWithCounter(dst io.Writer, src io.Reader, counter *uint64) (written int64, err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = io.ErrShortWrite
				}
			}
			written += int64(nw)
			atomic.AddUint64(counter, uint64(nw))
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// GetInstanceStats returns runtime stats for a rule.
func (m *Manager) GetInstanceStats(ruleID uint64) (activeConns int64, rxBytes uint64, txBytes uint64, isRunning bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inst, exists := m.instances[ruleID]
	if !exists {
		return 0, 0, 0, false
	}
	return atomic.LoadInt64(&inst.ActiveConns), atomic.LoadUint64(&inst.RxBytes), atomic.LoadUint64(&inst.TxBytes), true
}

// InitAll starts all enabled rules on application launch.
func InitAll() {
	db := model.UseDB()
	if db == nil {
		return
	}

	var rules []model.PortForwardRule
	if err := db.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return
	}

	for i := range rules {
		_ = GlobalManager.StartRule(&rules[i])
	}
}
