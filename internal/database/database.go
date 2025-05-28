package database

import (
	"database/sql"
	"fmt"
	"reminder-bot/internal/models"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

func (d DatabaseConfig) ToString() string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host,
		d.Port,
		d.Username,
		d.Password,
		d.DBName,
	)
	return dsn
}

type Database struct {
	db *sql.DB
}

func New(config DatabaseConfig) (*Database, error) {
	db, err := sql.Open("postgres", config.ToString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateUser(userName string, chatID int) error {
	_, err := d.db.Exec("insert into users (username, chat_id) values ($1, $2)", userName, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) CreateReminder(userID int, content string, interval time.Duration) error {
	_, err := d.db.Exec(
		"insert into reminders (user_id, content, interval) values ($1, $2, make_interval(mins => $3))",
		userID,
		content,
		interval.Minutes(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (d *Database) UpdateInterval(id int64, interval time.Duration) error {
	_, err := d.db.Exec("update reminders set interval = make_interval(mins => $1) where id=$2", interval.Minutes(), id)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) UpdateContent(id int64, content string) error {
	_, err := d.db.Exec("update reminders set content = $1 where id = $2", content, id)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteReminder(id int64) error {
	_, err := d.db.Exec("delete from reminders where id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) UpdateActive(id int64, isActive bool) error {
	_, err := d.db.Exec("update reminders set is_active = $1 where id = $2", isActive, id)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetUser(userName string, chatID int) (models.User, error) {
	row := d.db.QueryRow("select id, username, chat_id from users where username = $1 and chat_id = $2", userName, chatID)
	err := row.Err()
	if err != nil {
		return models.User{}, err
	}

	var result models.User

	err = row.Scan(&result.ID, &result.UserName, &result.ChatID)
	if err != nil {
		return models.User{}, err
	}

	return result, nil
}

func (d *Database) GetReminders(isActive bool) ([]models.Reminder, error) {
	rows, err := d.db.Query("select id, user_id, content, (extract(epoch from interval) / 60)::int as \"interval_minutes\", is_active, last_checked from reminders where is_active = $1", isActive)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []models.Reminder

	for rows.Next() {
		var reminder models.Reminder

		var minutes int
		err = rows.Scan(&reminder.ID, &reminder.UserID, &reminder.Content, &minutes, &reminder.IsActive, &reminder.LastChecked)
		if err != nil {
			return nil, err
		}

		reminder.Interval = time.Minute * time.Duration(minutes)

		result = append(result, reminder)
	}

	return result, nil
}

func (d *Database) GetChatID(id int64) (int64, error) {
	row := d.db.QueryRow("select chat_id from users where id = $1", id)
	err := row.Err()
	if err != nil {
		return 0, err
	}

	var result int64
	err = row.Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (d *Database) UpdateLastCheched(id int64) error {
	_, err := d.db.Exec("update reminders set last_checked = $1 where id = $2", time.Now().UTC(), id)
	if err != nil {
		return err
	}

	return nil
}
