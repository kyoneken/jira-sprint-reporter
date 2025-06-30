# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
Jira Sprint Reporter - A Go-based command-line tool that fetches active sprints from Jira, allows interactive sprint selection, and displays stories/tasks in a reporting-friendly tab-separated format optimized for Confluence table pasting.

## Common Commands

```bash
# Build the application
go build -o jira-sprint-reporter

# Run the application
./jira-sprint-reporter

# Install dependencies
go mod tidy

# Format code
go fmt ./...

# Run tests (when implemented)
go test ./...
```

## Architecture

### Main Components
- **JiraClient**: HTTP client for Jira REST API interactions with basic authentication
- **Sprint/Issue Models**: Data structures for Jira API responses
- **Interactive Selection**: Uses promptui for terminal-based sprint selection
- **Environment Management**: Uses godotenv for .env file loading

### API Endpoints Used
- `/rest/agile/1.0/sprint?state=active` - Fetch active sprints
- `/rest/agile/1.0/sprint/{sprintId}/issue` - Fetch issues in a sprint

### Configuration
Environment variables are managed via .env file:
- `JIRA_URL`: Base URL for Jira instance
- `JIRA_EMAIL`: Email for basic authentication
- `JIRA_API_TOKEN`: API token for authentication

### Dependencies
- `github.com/joho/godotenv`: Environment variable management
- `github.com/manifoldco/promptui`: Interactive terminal prompts