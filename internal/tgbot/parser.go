package tgbot

import (
	"errors"
	"fmt"
	"reminder-bot/internal/models"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func parse(message *tgbotapi.Message) (*models.ReminderCommand, error) {

	userInput := message.Text
	var parts []string

	if strings.Contains(userInput, "\n") {
		parts = strings.Split(userInput, "\n")
	} else {
		parts = strings.Fields(userInput)
	}

	command := &models.ReminderCommand{Command: strings.TrimSpace(parts[0])}

	switch command.Command {
	case "create", "Create":
		if len(parts) < 3 {
			return nil, fmt.Errorf(fmt.Sprintf("ошибка create, len(parts):%d", len(parts)))
		}

		time, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга времени, parts[1]:%s", parts[1])
		}

		command.Time = time
		command.Text = strings.Join(parts[2:], " ")

	case "edit_time", "Edit_time", "edittime", "Edittime":
		if len(parts) != 3 {
			return nil, fmt.Errorf("ошибка edit_time, len(parts):%d", len(parts))
		}

		reminderId, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга айди, parts[1]:%s", parts[1])
		}

		reminderTime, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга времени, parts[2]:%s", parts[2])
		}

		command.Id = reminderId
		command.Time = reminderTime

	case "edit_text", "Edit_text", "edittext", "Edittext":
		if len(parts) < 3 {
			return nil, fmt.Errorf("ошибка edit_text, len(parts):%d", len(parts))
		}

		reminderId, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга айди, parts[1]:%s", parts[1])
		}

		command.Id = reminderId
		command.Text = strings.Join(parts[2:], " ")

	case "toggle", "Toggle":
		if len(parts) < 3 {
			return nil, fmt.Errorf("ошибка toggle, len(parts):%d", len(parts))
		}

		reminderId, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга айди, parts[1]:%s", parts[1])
		}

		command.Id = reminderId
		command.Text = parts[2]

	case "delete", "Delete":
		if len(parts) != 2 {
			return nil, fmt.Errorf("ошибка delete/toggle, len(parts):%d", len(parts))
		}

		reminderId, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга айди, parts[1]:%s", parts[1])
		}

		command.Id = reminderId

	case "list", "List":
		if len(parts) > 1 {
			return nil, fmt.Errorf("ошибка list, len(parts):%d", len(parts))
		}

	default:
		return nil, errors.New("не нашёл условия")
	}

	return command, nil
}
