package dto

type GenerateImageTask struct {
	ID            string `json:"ID"`
	Prompt        string `json:"prompt"`
	AuthorId      string `json:"AuthorId"`
	ArticleDateId string `json:"ArticleDateId"`
}
