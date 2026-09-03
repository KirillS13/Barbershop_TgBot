package texts

import (
	"fmt"
	"time"
)

type DateOption struct {
	Label string
	Value string
}

func GetNextDays(count int) []DateOption {
	var days []DateOption
	months := []string{"янв", "фев", "мар", "апр", "май", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"}
	weekdays := []string{"Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

	now := time.Now()

	for i := 1; i <= count; i++ {
		date := now.AddDate(0, 0, i)
		weekday := weekdays[date.Weekday()]
		month := months[date.Month()-1]

		label := fmt.Sprintf("%s, %d %s", weekday, date.Day(), month)
		days = append(days, DateOption{Label: label, Value: date.Format("2006-01-02")})
	}
	return days
}
