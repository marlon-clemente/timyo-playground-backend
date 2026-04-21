package handlers


type FormCreateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=255"`
}

type SaveFormVersionRequest	 struct {
	Props any `json:"props" validate:"required"`
}