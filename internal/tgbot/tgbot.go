package tgbot

import (
	"fmt"
	"log"
	"reminder-bot/internal/database"
	"reminder-bot/internal/models"
	"strconv"
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

func getUsersReminders(db *database.Database, userID int, state bool) ([]models.Reminder, error) {
	reminders, err := db.GetReminders(state)
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
			reminderCommand, err := parse(update.Message)

			if err != nil {
				log.Printf(err.Error())
				continue
			}

			user, err := registerUser(db, update.Message.From.UserName, int(update.Message.Chat.ID))
			if err != nil {
				log.Printf(err.Error())
				continue
			}

			reminders, err_ := getUsersReminders(db, user.ID, true)

			switch reminderCommand.Command {
			case "create", "Create":
				err := db.CreateReminder(user.ID, reminderCommand.Text, time.Duration(reminderCommand.Time)*time.Duration(time.Minute)) // не возвращает ID поэтому невозможно указать ID при отправки сообщения
				if err != nil {
					log.Printf(err.Error())
					continue
				}

				SendMessage(fmt.Sprintf("Напоминание будет приходить каждые %s минут с текстом:\n%s", time.Duration(reminderCommand.Time)*time.Duration(time.Minute), reminderCommand.Text),
					int64(user.ChatID),
					bot)

			case "edit_time", "Edit_time", "edittime", "Edittime":
				err := db.UpdateInterval(int64(reminderCommand.Id), time.Duration(reminderCommand.Time)*time.Duration(time.Minute))
				if err != nil {
					log.Printf(err.Error())

					SendMessage(fmt.Sprintf("Не удалось изменить время для напоминания \n №%d - %s", reminderCommand.Id, reminderCommand.Text),
						int64(user.ChatID),
						bot)

					continue
				}

				SendMessage(fmt.Sprintf("Теперь напоминание №%d будет приходить каждые %s минут", reminderCommand.Id, time.Duration(reminderCommand.Time)*time.Duration(time.Minute)),
					int64(user.ChatID),
					bot)

			case "edit_text", "Edit_text", "edittext", "Edittext":
				err := db.UpdateContent(int64(reminderCommand.Id), reminderCommand.Text)
				if err != nil {
					SendMessage(fmt.Sprintf("Не удалось изменить текст для напоминания \n №%d - %s", reminderCommand.Id, reminderCommand.Text),
						int64(user.ChatID),
						bot)

					continue
				}

				SendMessage(fmt.Sprintf("Теперь напоминание №%d будет приходить с текстом:\n%s", reminderCommand.Id, reminderCommand.Text),
					int64(user.ChatID),
					bot)

			case "delete", "Delete":
				err := db.DeleteReminder(int64(reminderCommand.Id))
				if err != nil {
					SendMessage(fmt.Sprintf("Не удалось удалить напоминание №%d", reminderCommand.Id),
						int64(user.ChatID),
						bot)

					continue
				}

				SendMessage(fmt.Sprintf("Напоминание №%d удалено", reminderCommand.Id),
					int64(user.ChatID),
					bot)

			case "toggle", "Toggle":
				state, err := strconv.ParseBool(reminderCommand.Text)
				if err != nil {
					SendMessage(fmt.Sprintf("Не удалось обработать состояние \"%s\"", reminderCommand.Id, reminderCommand.Text),
						int64(user.ChatID),
						bot)

					continue
				}

				err = db.UpdateActive(int64(reminderCommand.Id), state)
				if err != nil {
					log.Printf(err.Error())
					SendMessage(fmt.Sprintf("Не удалось сменить состояние напоминания №%d", reminderCommand.Id),
						int64(user.ChatID),
						bot)

					continue
				}

				SendMessage(fmt.Sprintf("Состояние напоминания №%d изменено на %s", reminderCommand.Id, reminderCommand.Text),
					int64(user.ChatID),
					bot)

			case "list", "List":
				var reminders_ []models.Reminder
				reminders_, err_ = getUsersReminders(db, user.ID, false)
				if err_ != nil {
					log.Printf(err.Error())
					SendMessage(fmt.Sprintf("Не удалось создать список напоминаний"),
						int64(user.ChatID),
						bot)
					continue
				}

				combined := make([]models.Reminder, len(reminders)+len(reminders_))

				copy(combined, reminders)
				copy(combined[len(reminders):], reminders_)

				SendList(combined,
					int64(user.ChatID),
					bot)
			}
		}
	}
}
