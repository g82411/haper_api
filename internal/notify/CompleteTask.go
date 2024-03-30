package notify

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"hyper_api/internal/config"
	"hyper_api/internal/dto"
	"strings"
)

type CompleteTaskMessage struct {
	TaskID      string `json:"taskID"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	ArticleId   string `json:"articleID"`
	RecipientID string `json:"recipientID"`
}

type CompleteTaskAction struct {
	Action string `json:"action"`
	Body   string `json:"body"`
}

func CompleteTask(task dto.GenerateImageTask) error {
	setting := config.GetConfig()
	connection, _, err := websocket.DefaultDialer.Dial(setting.WebSocketURL, nil)
	if err != nil {
		return fmt.Errorf("error dialing websocket %v", err)
	}
	defer connection.Close()
	articleIdSplit := strings.Split(task.ArticleDateId, "_")
	articleId := articleIdSplit[0]
	message := CompleteTaskMessage{
		TaskID:      task.ID,
		RecipientID: task.AuthorId,
		ArticleId:   articleId,
		Type:        "imageGenerate",
		Status:      "success",
	}

	stringMessage, err := json.Marshal(message)
	stringMessage, err = json.Marshal(CompleteTaskAction{
		Action: "taskDone",
		Body:   string(stringMessage),
	})
	if err != nil {
		return fmt.Errorf("error marshalling message %v", err)
	}
	err = connection.WriteMessage(websocket.TextMessage, stringMessage)
	if err != nil {
		return fmt.Errorf("error writing message %v", err)
	}
	return nil
}
