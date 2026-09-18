package ddns

import (
	"sync"
	"time"

	"github.com/0xJacky/Nginx-UI/model"
	"github.com/uozi-tech/cosy/logger"
)

type Scheduler struct {
	stopChan chan struct{}
	mu       sync.Mutex
	running  bool
}

var GlobalScheduler = &Scheduler{
	stopChan: make(chan struct{}),
}

// Start runs the periodic DDNS check worker.
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.stopChan:
				return
			case <-ticker.C:
				s.checkAllTasks()
			}
		}
	}()
}

func (s *Scheduler) checkAllTasks() {
	db := model.UseDB()
	if db == nil {
		return
	}

	var tasks []model.DDNSTask
	if err := db.Where("enabled = ?", true).Find(&tasks).Error; err != nil {
		return
	}

	now := time.Now()
	for i := range tasks {
		task := &tasks[i]
		interval := task.IntervalSeconds
		if interval <= 0 {
			interval = 300
		}

		if task.LastRunAt != nil && now.Sub(*task.LastRunAt) < time.Duration(interval)*time.Second {
			continue
		}

		go RunTask(task.ID)
	}
}

// RunTask executes IP check and DNS update for a specific task.
func RunTask(taskID uint64) error {
	db := model.UseDB()
	var task model.DDNSTask
	if err := db.First(&task, taskID).Error; err != nil {
		return err
	}

	now := time.Now()
	task.LastRunAt = &now

	// 1. Fetch current IP
	currentIP, err := GetIPByMethod(task.IPType, task.IPMethod, task.IPInterface, task.IPUrl)
	if err != nil {
		task.LastStatus = "failed"
		task.LastError = err.Error()
		_ = db.Save(&task)
		logger.Warnf("[DDNS Task %s] Failed to retrieve IP: %v", task.Name, err)
		return err
	}

	// 2. Check if IP changed
	lastIP := task.LastIPv4
	if task.IPType == "ipv6" {
		lastIP = task.LastIPv6
	}

	if currentIP == lastIP && task.LastStatus == "success" {
		// IP unchanged
		_ = db.Save(&task)
		return nil
	}

	// 3. Update DNS Record
	err = UpdateRecord(&task, currentIP)
	if err != nil {
		task.LastStatus = "failed"
		task.LastError = err.Error()
		_ = db.Save(&task)
		logger.Errorf("[DDNS Task %s] DNS update error: %v", task.Name, err)
		return err
	}

	// 4. Update state
	if task.IPType == "ipv6" {
		task.LastIPv6 = currentIP
	} else {
		task.LastIPv4 = currentIP
	}
	task.LastStatus = "success"
	task.LastError = ""
	_ = db.Save(&task)

	logger.Infof("[DDNS Task %s] IP successfully updated to %s", task.Name, currentIP)
	return nil
}
