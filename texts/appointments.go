package texts

import "my-tg-bot/models"

func GetAllAppointments(appointments []models.Appointment) string {
	allAppointments := "📋 Ваши записи \n"
	for _, appointment := range appointments {
		appointmentDate := appointment.Time.Format("02.01.2006")
		hour := appointment.Time.Format("15:00")
		var master string
		var service string
		var status string
		switch appointment.MasterId {
		case 1:
			master = "Алекс"
		case 2:
			master = "Дмитрий"
		}
		switch appointment.ServerId {
		case 1:
			service = "Стрижка"
		case 2:
			service = "Борода"
		case 3:
			service = "Комбо"
		}
		switch appointment.Status {
		case "confirmed":
			status = "🟢 Подтверждено"
		case "rejected":
			status = "Отклонено"
		}
		app := "🗓" + appointmentDate + " - " + hour + "\n ✂️ Услуга: " + service + "\n👤 Мастер: " + master + "\n Статус: " + status + "\n\n"
		allAppointments += app
	}

	return allAppointments
}
