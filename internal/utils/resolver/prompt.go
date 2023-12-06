package resolver

import (
	"fmt"
	"hyper_api/internal/dto"
)

func GenerateImagePrompt(req dto.GenerateImageRequest) string {
	var ret string
	styles := []string{
		"單一線條、色塊",
		"平面插畫",
		"擬真",
	}
	style := styles[req.Style]
	if req.Action == 1 {
		template := "請繪製%v，圖案以%v方式呈現，不要太抽象。"
		ret = fmt.Sprintf(template, req.Items[0], style)
		return ret
	}
	if req.Action == 2 {
		template := "請繪製%v，圖案以%v方式呈現，人要有五官、不要太抽象、不要太多裝飾物。"
		ret = fmt.Sprintf(template, req.Items[0], style)
		return ret
	}
	relationTemplate := "%v在%v的%v"
	relation := fmt.Sprintf(relationTemplate, req.Items[0], req.Items[1], req.Relation)
	template := "請繪製%v，圖案以%v方式呈現，圖案以擬真方式呈現，請以台灣或亞洲場景為主。"
	ret = fmt.Sprintf(template, relation, style)
	return ret
}
