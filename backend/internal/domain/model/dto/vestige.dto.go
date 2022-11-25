package dto

type VestigeDto struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Parent  string `json:"parent" validate:"ethAddr" transform:"hex"`
	Head    string `json:"head" validate:"ethAddr" transform:"hex"`
}
