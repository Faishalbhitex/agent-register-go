```markdown
# Agent Register Service

A lightweight microservice for registering and discovering A2A (Agent-to-Agent) protocol compatible agents. This service acts as a registry/directory that allows agent clients to discover available agents based on their capabilities and skills.

## Features

- **Auto-discovery**: Automatically fetch agent cards from A2A agent servers
- **REST API**: Simple CRUD operations for agent management
- **Skills-based discovery**: Query agents by their capabilities and skills
- **Multi-platform support**: Works with A2A agents written in TypeScript, Python, Java, Go, .NET
- **SQLite storage**: Lightweight, serverless database for persistence
- **Standards compliant**: Uses official Google A2A SDK for agent card resolution

## Quick Start

### Prerequisites

- Go 1.25+ 
- SQLite3 (for database inspection, optional)

### Installation

```bash
# Clone or download the project
cd agent-register-go

# Install dependencies
go mod tidy

# Run the service
go run main.go
```

The service will start on `http://localhost:8080`

### Basic Usage

#### 1. Register an Agent

```bash
curl -X POST http://localhost:8080/agents \
  -H "Content-Type: application/json" \
  -d '{"url": "http://your-agent-server:4000"}'
```

**Response:**
```json
{
  "message": "Agent registered successfully",
  "agent": {
    "id": 1,
    "name": "Weather Agent",
    "skills": ["weather", "forecast", "location-specific"],
    "description": "Provides weather information for specific regions",
    "url": "http://your-agent-server:4000",
    "created_at": "2025-09-21 19:33:06"
  }
}
```

#### 2. Discover All Agents

```bash
curl http://localhost:8080/agents
```

**Response:**
```json
{
  "count": 1,
  "agents": [
    {
      "id": 1,
      "name": "Weather Agent", 
      "skills": ["weather", "forecast", "location-specific"],
      "description": "Provides weather information for specific regions",
      "url": "http://your-agent-server:4000",
      "created_at": "2025-09-21 19:33:06"
    }
  ]
}
```

#### 3. Get Specific Agent

```bash
curl http://localhost:8080/agents/1
```

#### 4. Delete Agent

```bash
curl -X DELETE http://localhost:8080/agents/1
```

#### 5. Health Check

```bash
curl http://localhost:8080/health
```

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/agents` | Register new agent by URL |
| GET    | `/agents` | List all registered agents |
| GET    | `/agents/:id` | Get agent by ID |
| DELETE | `/agents/:id` | Delete agent by ID |
| GET    | `/health` | Service health check |

### Error Responses

- `400 Bad Request` - Invalid request format or agent URL unreachable
- `404 Not Found` - Agent not found
- `409 Conflict` - Agent with URL already registered
- `500 Internal Server Error` - Database or server error

## How It Works

1. **Agent Registration**: When you POST an agent URL, the service:
   - Fetches the agent card from `{url}/.well-known/agent-card.json`
   - Extracts agent name, skills, and description
   - Stores the information in SQLite database
   - Returns success/error response

2. **Skills Extraction**: Skills are automatically extracted from:
   - Agent card skill names
   - Agent card skill tags
   - Fallback to "general" if no skills found

3. **Discovery**: Clients can query the registry to find agents matching their requirements

## Database Schema

The service uses SQLite with the following schema:

```sql
CREATE TABLE agents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    skills TEXT NOT NULL,           -- JSON array of skills
    description TEXT,
    url TEXT NOT NULL UNIQUE,       -- Agent server URL
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Database Operations

```bash
# Connect to database
sqlite3 agents.db

# View all agents
.mode column
SELECT * FROM agents;

# View schema
.schema

# Count agents
SELECT COUNT(*) FROM agents;
```

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Agent Server  │    │  Agent Register  │    │  Agent Client   │
│  (TS/Py/Go/etc) │────│   (This Service) │────│   (Discovery)   │
│                 │    │                  │    │                 │
│ Exposes A2A     │    │ Stores & Serves  │    │ Queries & Uses  │
│ Agent Card      │    │ Agent Directory  │    │ Agent Info      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Configuration

The service uses sensible defaults and requires minimal configuration:

- **Port**: 8080 (hardcoded in main.go)
- **Database**: `./agents.db` (SQLite file)
- **Timeout**: 30 seconds for agent card fetching

## Development

### Project Structure

```
agent-register-go/
├── main.go              # Entry point and server setup
├── database/
│   └── database.go      # SQLite connection and schema
├── models/
│   └── agent_card.go    # Agent model and business logic
├── handlers/
│   └── agent_handler.go # HTTP request handlers
├── routers/
│   └── router.go        # Route definitions
├── agents.db           # SQLite database (auto-created)
├── go.mod              # Go module definition
└── README.md           # This file
```

### Adding New Features

1. **Model changes**: Update `models/agent_card.go`
2. **API endpoints**: Add handlers in `handlers/` and routes in `routers/`
3. **Database changes**: Update schema in `database/database.go`

## Production Deployment

### Basic Production Setup

```bash
# Set to production mode
export GIN_MODE=release

# Run with custom port
PORT=8080 go run main.go
```

### Recommended Production Considerations

- Use environment variables for configuration
- Add authentication/authorization if needed
- Set up reverse proxy (nginx/Apache)
- Configure proper logging
- Set up monitoring and health checks
- Use proper database backup strategy

## Troubleshooting

### Common Issues

**Agent registration fails:**
- Ensure the agent server is running and accessible
- Check if agent exposes `.well-known/agent-card.json` endpoint
- Verify agent card format matches A2A specification

**Database errors:**
- Check file permissions for `agents.db`
- Ensure SQLite3 is available
- Verify disk space availability

**Connection issues:**
- Check firewall settings
- Verify port 8080 is available
- Test agent server connectivity manually

### Debugging

```bash
# Enable debug logging
export GIN_MODE=debug

# Check agent card manually
curl http://your-agent-server/.well-known/agent-card.json

# Test connectivity
curl http://localhost:8080/health
```

## Contributing

This is a microservice component of a larger A2A ecosystem project. For improvements:

1. Focus on registry functionality
2. Maintain API compatibility  
3. Keep dependencies minimal
4. Ensure SQLite compatibility

## License

This project is part of an A2A ecosystem implementation following the A2A protocol specification.

## Related Projects

- **A2A Protocol**: https://github.com/a2aproject/A2A
- **A2A Go SDK**: https://github.com/a2aproject/a2a-go
- **MCP Protocol**: https://github.com/modelcontextprotocol
```
