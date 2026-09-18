package logcenter

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/0xJacky/Nginx-UI/internal/docker_mgr"
	"github.com/0xJacky/Nginx-UI/settings"
)

type LogSource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // "file", "docker", "system", "audit"
	Description string `json:"description"`
}

// GetSources returns all queryable log sources.
func GetSources() []LogSource {
	sources := []LogSource{
		{
			ID:          "nginx_access",
			Name:        "Nginx 访问日志 (Access Log)",
			Type:        "file",
			Description: "记录所有进入 Nginx 的 HTTP/HTTPS 客户端请求及状态码",
		},
		{
			ID:          "nginx_error",
			Name:        "Nginx 错误日志 (Error Log)",
			Type:        "file",
			Description: "记录 Nginx 进程错误、配置报错、上游断开等异常信息",
		},
		{
			ID:          "system",
			Name:        "系统日志 (System Log / Journal)",
			Type:        "system",
			Description: "记录操作系统内核、守护进程及系统服务的系统级日志",
		},
		{
			ID:          "audit",
			Name:        "面板操作审计日志",
			Type:        "audit",
			Description: "记录管理员在控制面板进行的所有配置修改、操作与触发记录",
		},
	}

	// Append running Docker containers as selectable log sources
	containers, err := docker_mgr.ListContainers()
	if err == nil {
		for _, c := range containers {
			name := c.ID
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}
			sources = append(sources, LogSource{
				ID:          "docker_" + c.ID,
				Name:        "Docker: " + name,
				Type:        "docker",
				Description: "容器 " + name + " (" + c.Image + ") 的标准控制台输出",
			})
		}
	}

	return sources
}

// QueryLogs searches logs from the selected source.
func QueryLogs(source string, keyword string, limit int) ([]string, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}

	if strings.HasPrefix(source, "docker_") {
		containerID := strings.TrimPrefix(source, "docker_")
		logs, err := docker_mgr.GetContainerLogs(containerID, "300")
		if err != nil {
			return []string{"Failed to read docker logs: " + err.Error()}, nil
		}
		return filterLines(strings.Split(logs, "\n"), keyword, limit), nil
	}

	switch source {
	case "nginx_access":
		path := settings.NginxSettings.AccessLogPath
		if path == "" {
			path = "/var/log/nginx/access.log"
		}
		return readTailFile(path, keyword, limit)
	case "nginx_error":
		path := settings.NginxSettings.ErrorLogPath
		if path == "" {
			path = "/var/log/nginx/error.log"
		}
		return readTailFile(path, keyword, limit)
	case "system":
		if runtime.GOOS == "linux" {
			out, err := exec.Command("journalctl", "-n", "200", "--no-pager").CombinedOutput()
			if err == nil {
				return filterLines(strings.Split(string(out), "\n"), keyword, limit), nil
			}
		}
		// Fallback to syslog
		for _, p := range []string{"/var/log/syslog", "/var/log/messages", "/var/log/system.log"} {
			if _, err := os.Stat(p); err == nil {
				return readTailFile(p, keyword, limit)
			}
		}
		return []string{"System log is not accessible in this environment"}, nil
	default:
		return []string{"Select a log source from the list"}, nil
	}
}

func readTailFile(filePath string, keyword string, limit int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{"Log file not found or unreadable: " + filePath}, nil
	}
	defer file.Close()

	var allLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	return filterLines(allLines, keyword, limit), nil
}

func filterLines(lines []string, keyword string, limit int) []string {
	var filtered []string
	kw := strings.ToLower(keyword)

	// Traverse from bottom to top (most recent first)
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if kw == "" || strings.Contains(strings.ToLower(line), kw) {
			filtered = append(filtered, line)
			if len(filtered) >= limit {
				break
			}
		}
	}

	// Reverse back so oldest is top, newest is bottom
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}

	return filtered
}
