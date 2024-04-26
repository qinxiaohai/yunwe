package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// SendMessage sends a message through the Telegram Bot API
func SendMessage(text string) error {
	chatID := "-4184774560"
	url := "https://api.telegram.org/bot6399326213:AAEg44i3AXPj1_-xPZzm7S70-j1r5nf7tGw/sendMessage"

	// 构建消息体
	message := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 发送请求
	_, err = http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	return err
}
