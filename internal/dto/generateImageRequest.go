package dto

type GenerateImageRequest struct {
	Action   int
	Items    []string
	Relation string
	Style    int
}
