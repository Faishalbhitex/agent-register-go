// handlers/agent_handler.go
package handlers

import (
	"agent-register-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /agents - Register agent by URL
func RegisterAgent(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	// Validate URL format
	if request.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	// Fetch agent card and save to database
	agent, err := models.CreateAgentFromURL(request.URL)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: agents.url" {
			c.JSON(http.StatusConflict, gin.H{"error": "Agent with this URL already registered"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Agent registered successfully",
		"agent":   agent,
	})
}

// GET /agents - Get all agents with optional availability filter
func GetAllAgents(c *gin.Context) {
	// Check for availability filter
	availableOnly := c.Query("available") == "true"

	agents, err := models.GetAllAgents(availableOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch agents"})
		return
	}

	response := gin.H{
		"count":  len(agents),
		"agents": agents,
	}

	// Add filter info to response
	if availableOnly {
		response["filter"] = "available_only"
		response["note"] = "Showing only agents that are currently available or recently active"
	} else {
		response["filter"] = "all_registered"
		response["note"] = "Showing all registered agents regardless of availability"
	}

	c.JSON(http.StatusOK, response)
}

// GET /agents/:id - Get agent by ID
func GetAgentByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	agent, err := models.GetAgentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch agent"})
		return
	}

	if agent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// POST /agents/:id/heartbeat - Update agent heartbeat
func UpdateAgentHeartbeat(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	err = models.UpdateAgentHeartbeat(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update heartbeat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Agent heartbeat updated",
		"timestamp": "now",
		"status":    "available",
	})
}

// DELETE /agents/:id - Delete agent by ID
func DeleteAgent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	err = models.DeleteAgentByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete agent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}
