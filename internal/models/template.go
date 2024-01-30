package models

import "context"

type Template struct {
	ID     string
	Prompt string
	Title  string
}

func (template *Template) TableName(ctx context.Context) string {
	stage := ctx.Value("stage").(string)
	return stage + "_templates"
}

func (template *Template) Deserialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"id":     template.ID,
		"prompt": template.Prompt,
		"title":  template.Title,
	}, nil
}

func (template *Template) Serialize(av map[string]interface{}) {
	template.ID = av["id"].(string)
	template.Prompt = av["prompt"].(string)
	template.Title = av["title"].(string)
}
