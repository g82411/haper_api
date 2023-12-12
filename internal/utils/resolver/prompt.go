package resolver

import (
	"fmt"
	"hyper_api/internal/dto"
)

func generateTaiwanSnackPrompt(req dto.GenerateImageRequest) string {
	snackContainer, snackName := req.Items[0], req.Items[1]
	sauce := ""
	if req.Comment != "" {
		sauce = fmt.Sprintf("%v", req.Comment)
	}
	mainDescription := fmt.Sprintf("幫我繪製台灣小吃，一個在%v的%v", snackContainer, snackName)
	if sauce != "" {
		mainDescription = fmt.Sprintf("%v，%v", mainDescription, sauce)
	}
	additionalDescription := "圖片只有食物及配料，食物及顏色接近最終呈現的樣子。\n以卡通插畫的方式繪製，底圖為白色，線條乾淨俐落。"
	return fmt.Sprintf("%v\n%v", mainDescription, additionalDescription)
}

func generateTaiwanFestival(req dto.GenerateImageRequest) string {
	color, item, shape := req.Items[0], req.Items[1], req.Items[2]
	comment := req.Comment
	if comment != "" {
		comment = fmt.Sprintf("，同時%v。", comment)
	}
	mainDescription := fmt.Sprintf("請幫我繪製台灣節慶用品，一個%v的%v，形狀為%v%v\n", color, item, shape, comment)
	additionalDescription := "圖像貼近台灣用品真實的模樣，圖片僅有物品而已，不要有任何文字，旁邊不要有裝飾物。\n以卡通插畫的方式繪製，底圖為白色，線條乾淨俐落"
	return fmt.Sprintf("%v%v", mainDescription, additionalDescription)
}

func generateSport(req dto.GenerateImageRequest) string {
	sportName := req.Items[0]
	age := req.Comment
	if age != "" {
		age = fmt.Sprintf("，人物年齡設定在%v時期。", age)
	}
	mainDescription := fmt.Sprintf("請幫我繪製運動%v%v\n", sportName, age)
	additionalDescription := "圖片貼近此運動會做出的動作。圖片僅有動作及配件而已，人物一位即可，人物要有清楚的五官，旁邊不要有其他元素。\n以卡通插畫的方式繪製，底圖為白色，線條乾淨俐落"
	return fmt.Sprintf("%v%v", mainDescription, additionalDescription)
}

func GenerateImagePrompt(req dto.GenerateImageRequest) string {
	var ret string
	styles := []string{
		"卡通插畫",
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
		template := "請繪製%v，圖案以%v方式呈現，人要有清楚的五官、不要太抽象、不要太多裝飾物。"
		ret = fmt.Sprintf(template, req.Items[0], style)
		return ret
	}
	if req.Action == 4 {
		return generateTaiwanSnackPrompt(req)
	}
	if req.Action == 5 {
		return generateTaiwanFestival(req)
	}
	if req.Action == 6 {
		return generateSport(req)
	}
	relationTemplate := "%v在%v的%v"
	relation := fmt.Sprintf(relationTemplate, req.Items[0], req.Items[1], req.Relation)
	template := "請繪製%v，圖案以%v方式呈現，請以台灣或亞洲場景為主。"
	ret = fmt.Sprintf(template, relation, style)
	return ret
}
