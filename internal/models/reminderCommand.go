package models

type ReminderCommand struct {
	Command string
	Id      int
	Time    int
	Text    string
}
