package notify

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"hyper_api/internal/config"
	"hyper_api/internal/dto"
	"net/url"
)

type CompleteTaskMessage struct {
	TaskID      string `json:"taskID"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	RecipientID string `json:"recipientID"`
}

func CompleteTask(task dto.GenerateImageTask) error {
	setting := config.GetConfig()
	u := url.URL{Scheme: "wss", Host: setting.WebSocketHost, Path: setting.WebSocketPath}

	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("error dialing websocket %v", err)
	}
	defer connection.Close()
	message := CompleteTaskMessage{
		TaskID:      task.TaskID,
		RecipientID: task.AuthorId,
		Type:        "imageGenerate",
		Status:      "success",
	}
	stringMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message %v", err)
	}
	err = connection.WriteMessage(websocket.TextMessage, stringMessage)
	if err != nil {
		return fmt.Errorf("error writing message %v", err)
	}
	return nil
}
