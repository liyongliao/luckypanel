package docker_mgr

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/uozi-tech/cosy/logger"
)

type DockerStatus struct {
	Available      bool   `json:"available"`
	ServerVersion  string `json:"server_version"`
	Containers     int    `json:"containers"`
	ContainersRunning int `json:"containers_running"`
	ContainersPaused  int `json:"containers_paused"`
	ContainersStopped int `json:"containers_stopped"`
	Images         int    `json:"images"`
	Message        string `json:"message"`
}

type ContainerItem struct {
	ID      string   `json:"id"`
	Names   []string `json:"names"`
	Image   string   `json:"image"`
	State   string   `json:"state"`
	Status  string   `json:"status"`
	Ports   []string `json:"ports"`
	Created int64    `json:"created"`
}

type ImageItem struct {
	ID          string   `json:"id"`
	RepoTags    []string `json:"repo_tags"`
	Size        int64    `json:"size"`
	Created     int64    `json:"created"`
}

type ComposeStack struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Content    string    `json:"content"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string    `json:"status"`
}

func getClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = cli.Ping(ctx)
	if err != nil {
		_ = cli.Close()
		return nil, err
	}
	return cli, nil
}

// GetStatus checks if Docker is available and returns server info.
func GetStatus() DockerStatus {
	cli, err := getClient()
	if err != nil {
		return DockerStatus{
			Available: false,
			Message:   "Docker daemon is not reachable (socket not mounted or docker service stopped)",
		}
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := cli.Info(ctx)
	if err != nil {
		return DockerStatus{
			Available: false,
			Message:   err.Error(),
		}
	}

	return DockerStatus{
		Available:         true,
		ServerVersion:     info.ServerVersion,
		Containers:        info.Containers,
		ContainersRunning: info.ContainersRunning,
		ContainersPaused:  info.ContainersPaused,
		ContainersStopped: info.ContainersStopped,
		Images:            info.Images,
		Message:           "Docker is connected",
	}
}

// ListContainers returns list of containers.
func ListContainers() ([]ContainerItem, error) {
	cli, err := getClient()
	if err != nil {
		return nil, errors.New("docker is not reachable: " + err.Error())
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	res := make([]ContainerItem, len(containers))
	for i, c := range containers {
		var ports []string
		for _, p := range c.Ports {
			if p.PublicPort > 0 {
				ports = append(ports, fmt.Sprintf("%d->%d/%s", p.PublicPort, p.PrivatePort, p.Type))
			} else {
				ports = append(ports, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
			}
		}

		res[i] = ContainerItem{
			ID:      c.ID[:12],
			Names:   c.Names,
			Image:   c.Image,
			State:   c.State,
			Status:  c.Status,
			Ports:   ports,
			Created: c.Created,
		}
	}
	return res, nil
}

// ContainerAction performs start, stop, restart, kill, remove on a container.
func ContainerAction(containerID, action string) error {
	cli, err := getClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch action {
	case "start":
		return cli.ContainerStart(ctx, containerID, container.StartOptions{})
	case "stop":
		timeout := 10
		return cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
	case "restart":
		timeout := 10
		return cli.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout})
	case "kill":
		return cli.ContainerKill(ctx, containerID, "SIGKILL")
	case "remove":
		return cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}

// GetContainerLogs returns recent stdout/stderr output.
func GetContainerLogs(containerID string, tail string) (string, error) {
	cli, err := getClient()
	if err != nil {
		return "", err
	}
	defer cli.Close()

	if tail == "" {
		tail = "200"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
		Timestamps: true,
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, reader)
	return buf.String(), nil
}

// ListImages returns Docker images.
func ListImages() ([]ImageItem, error) {
	cli, err := getClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	images, err := cli.ImageList(ctx, image.ListOptions{All: false})
	if err != nil {
		return nil, err
	}

	res := make([]ImageItem, len(images))
	for i, img := range images {
		id := img.ID
		if strings.HasPrefix(id, "sha256:") && len(id) > 19 {
			id = id[7:19]
		}
		res[i] = ImageItem{
			ID:       id,
			RepoTags: img.RepoTags,
			Size:     img.Size,
			Created:  img.Created,
		}
	}
	return res, nil
}

const composeDir = "data/compose"

// ListComposeStacks returns all saved compose stacks.
func ListComposeStacks() ([]ComposeStack, error) {
	_ = os.MkdirAll(composeDir, 0755)

	entries, err := os.ReadDir(composeDir)
	if err != nil {
		return nil, err
	}

	var stacks []ComposeStack
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		stackPath := filepath.Join(composeDir, entry.Name(), "docker-compose.yml")
		content, err := os.ReadFile(stackPath)
		if err != nil {
			continue
		}
		fi, _ := entry.Info()
		updatedAt := time.Now()
		if fi != nil {
			updatedAt = fi.ModTime()
		}

		stacks = append(stacks, ComposeStack{
			Name:      entry.Name(),
			Path:      stackPath,
			Content:   string(content),
			UpdatedAt: updatedAt,
			Status:    "ready",
		})
	}

	return stacks, nil
}

// SaveComposeStack creates or updates a compose stack.
func SaveComposeStack(name, content string) error {
	if name == "" || strings.Contains(name, "/") || strings.Contains(name, "..") {
		return errors.New("invalid stack name")
	}

	dir := filepath.Join(composeDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	stackPath := filepath.Join(dir, "docker-compose.yml")
	return os.WriteFile(stackPath, []byte(content), 0644)
}

// ComposeAction runs up, down, or restart for a stack.
func ComposeAction(name, action string) (string, error) {
	dir := filepath.Join(composeDir, name)
	stackPath := filepath.Join(dir, "docker-compose.yml")
	if _, err := os.Stat(stackPath); err != nil {
		return "", errors.New("stack not found: " + name)
	}

	var args []string
	switch action {
	case "up":
		args = []string{"compose", "-f", stackPath, "up", "-d"}
	case "down":
		args = []string{"compose", "-f", stackPath, "down"}
	case "restart":
		args = []string{"compose", "-f", stackPath, "restart"}
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}

	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	logger.Infof("Compose [%s %s]: %s", name, action, string(out))
	if err != nil {
		return string(out), fmt.Errorf("%s: %v", strings.TrimSpace(string(out)), err)
	}
	return string(out), nil
}
