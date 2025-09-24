package models

import (
	"agent-register-go/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/a2aproject/a2a-go/a2a"
	"github.com/a2aproject/a2a-go/a2aclient/agentcard"
)

type AgentCard struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Skills      []string `json:"skills"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Status      string   `json:"status"`
	LastSeenAt  *string  `json:"last_seen_at,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
}

// Fetch agent card from URL and create in database
func CreateAgentFromURL(url string) (*AgentCard, error) {
	// Fetch agent card using Google A2A SDK
	resolver := &agentcard.Resolver{
		BaseURL: url,
	}

	ctx := context.Background()
	card, err := resolver.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent card: %v", err)
	}

	// Extract skills from agent card
	skills := extractSkillsFromCard(card)

	// Create AgentCard struct
	agent := &AgentCard{
		Name:        card.Name,
		Skills:      skills,
		Description: card.Description,
		URL:         url,
		Status:      "registered", // Default status
	}

	// Save to database
	if err := agent.save(); err != nil {
		return nil, err
	}

	return agent, nil
}

// Extract skills from A2A AgentCard
func extractSkillsFromCard(card *a2a.AgentCard) []string {
	var skills []string

	// Extract from Skills field
	for _, skill := range card.Skills {
		if skill.Name != "" {
			skills = append(skills, skill.Name)
		}
		// Also extract tags if available
		for _, tag := range skill.Tags {
			if tag != "" && !contains(skills, tag) {
				skills = append(skills, tag)
			}
		}
	}

	// If no skills found, use default based on name/description
	if len(skills) == 0 {
		skills = []string{"general"}
	}

	return skills
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Save agent to database
func (a *AgentCard) save() error {
	skillsJSON, _ := json.Marshal(a.Skills)

	query := `INSERT INTO agents (name, skills, description, url, status) VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, a.Name, string(skillsJSON), a.Description, a.URL, a.Status)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	a.ID = int(id)
	return nil
}

// Get all agents with optional availability filter
func GetAllAgents(availableOnly bool) ([]AgentCard, error) {
	var query string

	if availableOnly {
		// Filter for agents with recent activity (last 5 minutes) or status 'available'
		query = `SELECT id, name, skills, description, url, status, last_seen_at, created_at 
				FROM agents 
				WHERE (status = 'available' OR (last_seen_at IS NOT NULL AND datetime(last_seen_at) > datetime('now', '-5 minutes')))
				ORDER BY created_at DESC`
	} else {
		query = `SELECT id, name, skills, description, url, status, last_seen_at, created_at FROM agents ORDER BY created_at DESC`
	}

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []AgentCard
	for rows.Next() {
		var agent AgentCard
		var skillsJSON string
		var createdAt time.Time
		var lastSeenAt sql.NullTime

		err := rows.Scan(&agent.ID, &agent.Name, &skillsJSON, &agent.Description,
			&agent.URL, &agent.Status, &lastSeenAt, &createdAt)
		if err != nil {
			return nil, err
		}

		// Parse skills JSON
		json.Unmarshal([]byte(skillsJSON), &agent.Skills)
		agent.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		// Handle nullable last_seen_at
		if lastSeenAt.Valid {
			lastSeenStr := lastSeenAt.Time.Format("2006-01-02 15:04:05")
			agent.LastSeenAt = &lastSeenStr
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

// Get agent by ID
func GetAgentByID(id int) (*AgentCard, error) {
	query := `SELECT id, name, skills, description, url, status, last_seen_at, created_at FROM agents WHERE id = ?`

	var agent AgentCard
	var skillsJSON string
	var createdAt time.Time
	var lastSeenAt sql.NullTime

	err := database.DB.QueryRow(query, id).Scan(
		&agent.ID, &agent.Name, &skillsJSON, &agent.Description,
		&agent.URL, &agent.Status, &lastSeenAt, &createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Parse skills JSON
	json.Unmarshal([]byte(skillsJSON), &agent.Skills)
	agent.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	// Handle nullable last_seen_at
	if lastSeenAt.Valid {
		lastSeenStr := lastSeenAt.Time.Format("2006-01-02 15:04:05")
		agent.LastSeenAt = &lastSeenStr
	}

	return &agent, nil
}

// Update agent heartbeat
func UpdateAgentHeartbeat(id int) error {
	query := `UPDATE agents SET status = 'available', last_seen_at = datetime('now') WHERE id = ?`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Update agent heartbeat by URL
func UpdateAgentHeartbeatByURL(url string) error {
	query := `UPDATE agents SET status = 'available', last_seen_at = datetime('now') WHERE url = ?`
	result, err := database.DB.Exec(query, url)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete agent by ID
func DeleteAgentByID(id int) error {
	query := `DELETE FROM agents WHERE id = ?`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
