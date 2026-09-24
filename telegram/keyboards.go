package telegram

import (
	"fmt"
	"my-tg-bot/texts"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✂️ Выбрать услугу", "show_services"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍🎨 Наша команда", "show_employees"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Ваши записи", "show_appointments"),
		),
	)
}

func ServicesMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✂️ Стрижка — 250 MDL", "service_haircut"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧔 Борода — 150 MDL", "service_beard"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥 Комплекс — 350 MDL", "service_combo"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 В главное меню", "show_main_menu"),
		),
	)
}

func BarberMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍🎨 Алекс (Top)", "barber_alex"),
			tgbotapi.NewInlineKeyboardButtonData("✂️ Дмитрий", "barber_dmitrii"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 В главное меню", "show_main_menu"),
		),
	)
}

func AlexMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Записаться к Алексу", "book_direct_alex"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👈 Другой мастер", "show_employees"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 В меню", "show_main_menu"),
		),
	)
}

func DmitriiMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Записаться к Дмитрию", "book_direct_dmitrii"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👈 Другой мастер", "show_employees"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 В меню", "show_main_menu"),
		),
	)
}

func HaircutMenu(service string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍🎨 Алекс (Top)", fmt.Sprintf("book_%s_alex", service)),
			tgbotapi.NewInlineKeyboardButtonData("✂️ Дмитрий", fmt.Sprintf("book_%s_dmitrii", service)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 Любой свободный мастер", fmt.Sprintf("book_%s_any", service)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к услугам", "show_services"),
		),
	)
}

func SelectDateMenu(service, barber string) tgbotapi.InlineKeyboardMarkup {
	days := texts.GetNextDays(6)

	row1 := []tgbotapi.InlineKeyboardButton{}
	row2 := []tgbotapi.InlineKeyboardButton{}

	for i, d := range days {
		cbData := fmt.Sprintf("date_%s_%s_%s", service, barber, d.Value)
		btn := tgbotapi.NewInlineKeyboardButtonData(d.Label, cbData)
		if i < 3 {
			row1 = append(row1, btn)
		} else {
			row2 = append(row2, btn)
		}
	}

	backCb := fmt.Sprintf("service_%s", service)

	return tgbotapi.NewInlineKeyboardMarkup(
		row1,
		row2,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к выбору мастера", backCb),
		),
	)
}

func SelectTimeMenu(service, barber, date string, bookedTimes []string) tgbotapi.InlineKeyboardMarkup {
	allSlots := []string{"10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00", "18:00"}

	bookedMap := make(map[string]bool)
	for _, t := range bookedTimes {
		bookedMap[t] = true
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	for _, slot := range allSlots {
		var btn tgbotapi.InlineKeyboardButton

		if bookedMap[slot] {
			btnText := strikeThrough(fmt.Sprintf("%s", slot))
			btn = tgbotapi.NewInlineKeyboardButtonData(btnText, "booked")
		} else {
			cbData := fmt.Sprintf("time_%s_%s_%s_%s", service, barber, date, slot)
			btn = tgbotapi.NewInlineKeyboardButtonData(slot, cbData)
		}

		currentRow = append(currentRow, btn)

		if len(currentRow) == 3 {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	backCb := fmt.Sprintf("book_%s_%s", service, barber)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к выбору даты", backCb)))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func strikeThrough(s string) string {
	var result []rune
	for _, r := range s {
		result = append(result, r, '\u0336')
	}
	return string(result)
}

func AppointmentMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 В главное меню", "show_main_menu"),
		),
	)
}
