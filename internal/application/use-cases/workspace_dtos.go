package usecases

type WorkspaceCreateInput struct {
	AgentID string
	Name string
}

type WorkspaceOutput struct {
	ID string `json:"id"`
	Name string `json:"name"`
	AgentID string `json:"agentID"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}