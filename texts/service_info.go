package texts

import "fmt"

func HaircutService(serviceName string) string {
	var price, duration, title string

	switch serviceName {
	case "haircut":
		price = "250 MDL"
		duration = "45 мин"
		title = "Мужская стрижка"
	case "beard":
		price = "150 MDL"
		duration = "30 мин"
		title = "Оформление бороды"
	case "combo":
		price = "350 MDL"
		duration = "60 мин"
		title = "Комплекс (Стрижка + Борода)"
	default:
		price = "-"
		duration = "-"
		title = "Услуга"
	}
	return fmt.Sprintf(
		"✂️ <b>Выбрана услуга: %s</b>\n"+
			"├ 💵 <b>Стоимость:</b> %s\n"+
			"└ ⏱ <b>Длительность:</b> %s\n\n"+
			"<i>Выберите мастера, к которому хотите записаться:</i>",
		title, price, duration,
	)
}
