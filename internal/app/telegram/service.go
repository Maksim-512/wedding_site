package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"wedding_website/internal/app/models"
)

type TelegramService struct {
	botToken string
	chatID   string
	myLogger slog.Logger
}

func NewTelegramService(botToken, chatID string, myLogger *slog.Logger) *TelegramService {
	return &TelegramService{
		botToken: botToken,
		chatID:   chatID,
		myLogger: *myLogger,
	}
}

type Message struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func (s *TelegramService) SendMessage(text string) error {
	s.myLogger.Info("Начало SendMessage с текстом", slog.String("text", text))

	message := Message{
		ChatID:    s.chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		s.myLogger.Error("Ошибка Marshal", slog.String("error", err.Error()))
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	s.myLogger.Info("Ссылка для отправки", slog.String("url", url))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.myLogger.Error("Ошибка отправки", slog.String("error", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.myLogger.Error("Не верный статус код", slog.Int("Code", resp.StatusCode))
		return fmt.Errorf("telegram API error: %s", resp.Status)
	}

	return nil
}

func (s *TelegramService) SendRSVPNotification(rsvp *models.RSVP) error {
	attendance := "❌ *Не придет*"
	if rsvp.Attendance {
		attendance = "✅ *Придет*"
	}

	emoji := "🎉"
	if !rsvp.Attendance {
		emoji = "😢"
	}

	companionInfo := "Один/одна"
	if rsvp.Companion != "" {
		companionInfo = fmt.Sprintf("С спутниками: *%s*", rsvp.Companion)
	}

	message := fmt.Sprintf(
		"%s *Новый ответ на приглашение!* %s\n\n"+
			"👤 *Имя:* %s\n"+
			"📋 *Присутствие:* %s\n"+
			"👥 *Компания:* %s\n"+
			"🕐 *Время:* %s",
		emoji, emoji,
		rsvp.Name,
		attendance,
		companionInfo,
		rsvp.CreatedAt.Format("02.01.2006 15:04"),
	)

	return s.SendMessage(message)
}
