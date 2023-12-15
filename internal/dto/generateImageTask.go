package dto

type GenerateImageTask struct {
	TaskID    string `json:"task_id"`
	ArticleId string `json:"article_id"`
	Prompt    string `json:"prompt"`
	AuthorId  string `json:"author_id"`
}
