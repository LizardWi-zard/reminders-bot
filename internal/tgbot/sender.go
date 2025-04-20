package tgbot

import (
	"fmt"
	"log"
	"reminder-bot/internal/database"
	"reminder-bot/internal/models"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(message string, userID int64, bot *tgbotapi.BotAPI) {
	log.Printf("попал в отправку сообщения")

	msg := tgbotapi.NewMessage(userID, message)
	bot.Send(msg)

	log.Printf("sending message to reminders to %d", userID)
}

func SendRemind(db *database.Database, bot *tgbotapi.BotAPI) {
	for {
		reminders, err := db.GetReminders(true)
		if err != nil {
			log.Println(err)
			continue
		}

		now := time.Now().UTC()
		for _, reminder := range reminders {
			lastChecked := reminder.LastChecked.UTC()
			nextCheckTime := lastChecked.Add(reminder.Interval)

			if now.After(nextCheckTime) || now.Equal(nextCheckTime) {
				userID, err := db.GetChatID(int64(reminder.UserID))
				if err != nil {
					log.Println(err)
					continue
				}

				msg := tgbotapi.NewMessage(userID, reminder.Content)
				bot.Send(msg)

				log.Printf("sending reminder %d to %d", reminder.UserID, userID)

				err = db.UpdateLastCheched(int64(reminder.ID))
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

func SendList(reminders []models.Reminder, userID int64, bot *tgbotapi.BotAPI) {
	var builder strings.Builder
	for _, rereminder := range reminders {
		line := fmt.Sprintf("%d - %s каждые %s ", rereminder.ID, rereminder.Content, rereminder.Interval)
		builder.WriteString(line)
	}

	msg := tgbotapi.NewMessage(userID, builder.String())
	bot.Send(msg)

	log.Printf("sending list of reminders to %d", userID)
}
