package docker

import (
	"net/http"

	"github.com/0xJacky/Nginx-UI/internal/docker_mgr"
	"github.com/gin-gonic/gin"
)

// GetStatus returns the Docker daemon status.
func GetStatus(c *gin.Context) {
	status := docker_mgr.GetStatus()
	c.JSON(http.StatusOK, status)
}

// ListContainers returns list of containers.
func ListContainers(c *gin.Context) {
	list, err := docker_mgr.ListContainers()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "data": []any{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ContainerAction performs start, stop, restart, kill, or remove on a container.
func ContainerAction(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := docker_mgr.ContainerAction(id, body.Action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetLogs returns container logs.
func GetLogs(c *gin.Context) {
	id := c.Param("id")
	tail := c.DefaultQuery("tail", "200")

	logs, err := docker_mgr.GetContainerLogs(id, tail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// ListImages returns Docker images.
func ListImages(c *gin.Context) {
	images, err := docker_mgr.ListImages()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "data": []any{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": images})
}

// ListCompose returns all compose stacks.
func ListCompose(c *gin.Context) {
	stacks, err := docker_mgr.ListComposeStacks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stacks)
}

// SaveCompose creates or updates a compose stack.
func SaveCompose(c *gin.Context) {
	var body struct {
		Name    string `json:"name" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := docker_mgr.SaveComposeStack(body.Name, body.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// ComposeAction executes up/down/restart on a stack.
func ComposeAction(c *gin.Context) {
	name := c.Param("name")
	var body struct {
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := docker_mgr.ComposeAction(name, body.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "output": out})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "output": out})
}
