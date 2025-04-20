package tgbot

import (
	"fmt"
	"log"
	"reminder-bot/internal/database"
	"reminder-bot/internal/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LaunchTheBot(db *database.Database) {
	bot, err := tgbotapi.NewBotAPI("") //enter tg-bot token here
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	go checkUpdates(db, bot)
	SendRemind(db, bot)
}

func registerUser(db *database.Database, userName string, chatID int) (models.User, error) {
	result, err := db.GetUser(userName, chatID)
	if err == nil {
		return result, nil
	}

	err = db.CreateUser(userName, chatID)
	if err != nil {
		return models.User{}, err
	}

	result, err = db.GetUser(userName, chatID)
	if err == nil {
		return result, nil
	}

	return models.User{
		UserName: userName,
		ChatID:   chatID,
	}, nil
}

func createReminder(db *database.Database, userID int, content string, interval time.Duration) error {
	err := db.CreateReminder(userID, content, interval)
	if err != nil {
		return err
	}
	return nil
}

func getUsersReminders(db *database.Database, userID int) ([]models.Reminder, error) {
	reminders, err := db.GetReminders(true)
	if err != nil {
		return nil, err
	}

	var usersReminds []models.Reminder
	for _, reminder := range reminders {
		if reminder.UserID == userID {
			usersReminds = append(usersReminds, reminder)
		}
	}

	return usersReminds, nil
}

func checkUpdates(db *database.Database, bot *tgbotapi.BotAPI) {

	//bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			reminderCommand, reminderInterval, reminderText, err := parse(update.Message)

			if err != nil {
				log.Printf(err.Error())
				continue
			}

			user, err := registerUser(db, update.Message.From.UserName, int(update.Message.Chat.ID))
			if err != nil {
				log.Printf(err.Error())
				continue
			}

			switch reminderCommand {
			case "new":
				{
					log.Printf("adding reminder")
					err = createReminder(db, user.ID, reminderText, time.Duration(reminderInterval)*time.Duration(time.Minute))
					if err != nil {
						log.Printf(err.Error())
						continue
					}

					SendMessage(fmt.Sprintf("Напоминание будет приходить каждые %s с текстом\n\n%s", time.Duration(reminderInterval)*time.Duration(time.Minute), reminderText),
						int64(user.ChatID),
						bot)
				}

			case "list":
				{
					log.Printf("showing the list of reminders")
					reminders, err := getUsersReminders(db, user.ID)
					if err != nil {
						log.Printf(err.Error())
						continue
					}

					SendList(reminders,
						int64(user.ChatID),
						bot)
				}
			}

			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

			//bot.Send(msg)
		}
	}
}
