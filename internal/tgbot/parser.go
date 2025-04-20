package tgbot

import (
	"errors"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func parse(msg *tgbotapi.Message) (string, int, string, error) {
	text := strings.TrimSpace(msg.Text)
	parts := strings.Split(text, "\n")

	log.Printf(text)

	if len(parts) < 3 {
		return "", 0, "", errors.New("Неправильный формат ввода")
	}

	log.Printf(parts[1])

	reminderCommand := parts[0]
	reminderInterval, _ := strconv.Atoi(parts[1])
	reminderText := parts[2]

	log.Printf(strconv.Itoa(reminderInterval))

	return reminderCommand, reminderInterval, reminderText, nil
}
