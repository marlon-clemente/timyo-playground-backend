package usecases

type CreateAgentInput struct {
	MemberID string
	Name string
	AvatarURL string

	AccessToken string
	
}

type AgentOutput struct {
	ID string `json:"id"`
	MemberID string `json:"memberID"`
	Name string `json:"name"`
	AvatarURL string `json:"avatarURL"`
	CreatedAt string `json:"createdAt"`
}