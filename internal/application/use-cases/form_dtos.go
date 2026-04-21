package usecases


type CreateFormInput struct {
	WorkspaceID string
	AgentID   string
	Name        string
	Description string
}

type SaveFormVersionInput struct {
	FormID string
	Props any
}

type DeleteFormInput struct {
	WorkspaceID string
	FormID string
}

type FormOutput struct {
	FormID string `json:"formID"`
	FormVersionID string `json:"formVersionID"`
	FormName string `json:"formName"`
	FormDescription string `json:"formDescription"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	Props any `json:"props,omitempty"`
}

type ListFormsInput struct {
	WorkspaceID string
}

type ListFormsOutput struct {
	Forms []FormOutput `json:"data"`
	Total int `json:"total"`
}