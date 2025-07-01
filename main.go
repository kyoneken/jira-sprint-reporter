package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

type JiraClient struct {
	BaseURL   string
	Email     string
	APIToken  string
	Client    *http.Client
}

type Sprint struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	State           string `json:"state"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	OriginBoardID   int    `json:"originBoardId"`
}

type SprintResponse struct {
	Values []Sprint `json:"values"`
}

type Issue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary   string `json:"summary"`
		IssueType struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		Assignee struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Epic struct {
			Name string `json:"name"`
		} `json:"epic"`
		StoryPoints *float64 `json:"customfield_10016"`
	} `json:"fields"`
}

type IssueResponse struct {
	Issues []Issue `json:"issues"`
}

func NewJiraClient() *JiraClient {
	return &JiraClient{
		BaseURL:  os.Getenv("JIRA_URL"),
		Email:    os.Getenv("JIRA_EMAIL"),
		APIToken: os.Getenv("JIRA_API_TOKEN"),
		Client:   &http.Client{},
	}
}

func (jc *JiraClient) makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(jc.Email, jc.APIToken)
	req.Header.Set("Content-Type", "application/json")

	return jc.Client.Do(req)
}

func (jc *JiraClient) GetActiveSprints() ([]Sprint, error) {
	// Get all boards first
	boardURL := fmt.Sprintf("%s/rest/agile/1.0/board", jc.BaseURL)
	
	resp, err := jc.makeRequest(boardURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch boards: %s", resp.Status)
	}

	var boardResponse struct {
		Values []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"values"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&boardResponse); err != nil {
		return nil, err
	}

	fmt.Printf("Found %d boards\n", len(boardResponse.Values))
	for _, board := range boardResponse.Values {
		fmt.Printf("  Board ID: %d, Name: %s\n", board.ID, board.Name)
	}

	if len(boardResponse.Values) == 0 {
		return nil, fmt.Errorf("no boards found")
	}

	var allSprints []Sprint
	
	// Check all boards for active sprints
	for _, board := range boardResponse.Values {
		sprintURL := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?state=active", jc.BaseURL, board.ID)
		
		resp, err = jc.makeRequest(sprintURL)
		if err != nil {
			fmt.Printf("Error fetching sprints for board %d: %v\n", board.ID, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to fetch sprints for board %d: %s\n", board.ID, resp.Status)
			continue
		}

		var sprintResponse SprintResponse
		if err := json.NewDecoder(resp.Body).Decode(&sprintResponse); err != nil {
			fmt.Printf("Error decoding sprints for board %d: %v\n", board.ID, err)
			continue
		}

		fmt.Printf("Board %d (%s) has %d active sprints\n", board.ID, board.Name, len(sprintResponse.Values))
		allSprints = append(allSprints, sprintResponse.Values...)
	}

	return allSprints, nil
}

func (jc *JiraClient) getSprintsDirectly() ([]Sprint, error) {
	// Try to access boards directly with minimal permissions
	boardURL := fmt.Sprintf("%s/rest/agile/1.0/board", jc.BaseURL)
	
	resp, err := jc.makeRequest(boardURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch boards: %s", resp.Status)
	}

	var boardResponse struct {
		Values []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"values"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&boardResponse); err != nil {
		return nil, err
	}

	fmt.Printf("Found %d boards\n", len(boardResponse.Values))
	
	if len(boardResponse.Values) == 0 {
		return nil, fmt.Errorf("no boards found")
	}

	var allSprints []Sprint
	
	// Get sprints from all boards
	for _, board := range boardResponse.Values {
		sprintURL := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?state=active", jc.BaseURL, board.ID)
		
		resp, err = jc.makeRequest(sprintURL)
		if err != nil {
			fmt.Printf("Error fetching sprints for board %d: %v\n", board.ID, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to fetch sprints for board %d: %s\n", board.ID, resp.Status)
			continue
		}

		var sprintResponse SprintResponse
		if err := json.NewDecoder(resp.Body).Decode(&sprintResponse); err != nil {
			fmt.Printf("Error decoding sprints for board %d: %v\n", board.ID, err)
			continue
		}

		fmt.Printf("Board %d (%s) has %d active sprints\n", board.ID, board.Name, len(sprintResponse.Values))
		allSprints = append(allSprints, sprintResponse.Values...)
	}

	return allSprints, nil
}

func (jc *JiraClient) GetSprintIssues(sprintID int) ([]Issue, error) {
	url := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d/issue?fields=key,assignee,epic", jc.BaseURL, sprintID)
	
	resp, err := jc.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch sprint issues: %s", resp.Status)
	}

	var issueResponse IssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issueResponse); err != nil {
		return nil, err
	}

	return issueResponse.Issues, nil
}

func selectSprint(sprints []Sprint) (*Sprint, error) {
	if len(sprints) == 0 {
		return nil, fmt.Errorf("no active sprints found")
	}

	if len(sprints) == 1 {
		return &sprints[0], nil
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Name | cyan }} ({{ .State }})",
		Inactive: "  {{ .Name | white }} ({{ .State }})",
		Selected: "✓ {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select an active sprint",
		Items:     sprints,
		Templates: templates,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &sprints[index], nil
}

func displayIssues(issues []Issue, jiraBaseURL string) {
	if len(issues) == 0 {
		fmt.Println("No issues found in this sprint.")
		return
	}

	// Tab-separated header for Confluence table
	fmt.Printf("LINK\tEPIC\tASSIGNEE\n")

	for _, issue := range issues {
		// Assignee
		assignee := "Unassigned"
		if issue.Fields.Assignee.DisplayName != "" {
			assignee = issue.Fields.Assignee.DisplayName
		}

		// Link
		link := fmt.Sprintf("%s/browse/%s", jiraBaseURL, issue.Key)

		// Epic
		epic := "-"
		if issue.Fields.Epic.Name != "" {
			epic = issue.Fields.Epic.Name
		}

		// Tab-separated output
		fmt.Printf("%s\t%s\t%s\n",
			link,
			epic,
			assignee,
		)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	client := NewJiraClient()

	// 環境変数の確認とデバッグ出力
	fmt.Printf("JIRA_URL: %s\n", client.BaseURL)
	fmt.Printf("JIRA_EMAIL: %s\n", client.Email)
	if client.APIToken != "" {
		fmt.Printf("JIRA_API_TOKEN: ***SET***\n")
	} else {
		fmt.Printf("JIRA_API_TOKEN: ***NOT SET***\n")
	}

	if client.BaseURL == "" || client.Email == "" || client.APIToken == "" {
		log.Fatal("Please set JIRA_URL, JIRA_EMAIL, and JIRA_API_TOKEN environment variables")
	}

	fmt.Println("Fetching active sprints...")
	sprints, err := client.GetActiveSprints()
	if err != nil {
		log.Fatalf("Error fetching sprints: %v", err)
	}

	selectedSprint, err := selectSprint(sprints)
	if err != nil {
		log.Fatalf("Error selecting sprint: %v", err)
	}

	fmt.Printf("\nFetching issues for sprint: %s\n", selectedSprint.Name)
	issues, err := client.GetSprintIssues(selectedSprint.ID)
	if err != nil {
		log.Fatalf("Error fetching sprint issues: %v", err)
	}

	displayIssues(issues, client.BaseURL)
}