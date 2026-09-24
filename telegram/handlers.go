package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"my-tg-bot/models"
	"my-tg-bot/services"
	"my-tg-bot/texts"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot                *tgbotapi.BotAPI
	db                 *sql.DB
	clientService      *services.ClientService
	mu                 sync.RWMutex
	userState          map[int]*models.BookingState
	appointmentService *services.AppointmentService
}

func NewHandler(bot *tgbotapi.BotAPI, db *sql.DB, clientService *services.ClientService, appointmentService *services.AppointmentService) *Handler {
	return &Handler{bot: bot, db: db, clientService: clientService, userState: make(map[int]*models.BookingState), appointmentService: appointmentService}
}

func (h *Handler) HandleCallBack(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID

	clientID, err := h.getOrCreateClient(query.From)
	if err != nil {
		log.Printf("Error processing client from callback: %v", err)
	}

	_ = clientID

	defer func() {
		callBackConfig := tgbotapi.NewCallback(query.ID, "")
		h.bot.Request(callBackConfig)
	}()
	appointments, err := h.appointmentService.GetAppointments(query.From.ID)
	if err != nil {
		log.Printf("Error processing appointments: %v", err)
	}

	switch query.Data {
	case "show_main_menu":
		h.editMessage(chatID, messageID, texts.GetMainMenuText(), MainMenu())
		return
	case "show_services":
		h.editMessage(chatID, messageID, texts.GetPriceList(), ServicesMenu())
		return
	case "show_employees":
		h.editMessage(chatID, messageID, texts.GetAllBarbers(), BarberMenu())
		return
	case "show_appointments":
		h.editMessage(chatID, messageID, texts.GetAllAppointments(appointments), AppointmentMenu())
	}
	parts := strings.Split(query.Data, "_")
	action := parts[0]
	switch action {
	case "barber":
		h.handleBarber(chatID, messageID, parts[1])
	case "service":
		h.handleService(chatID, messageID, parts[1])
	case "book":
		if len(parts) >= 3 {
			barberID := parts[2]
			serviceID := parts[1]
			h.handleBarberChosen(chatID, messageID, barberID, serviceID, int(query.From.ID))
		}
	case "date":
		if len(parts) >= 4 {
			barberID := parts[2]
			serviceID := parts[1]
			date := parts[3]
			h.handleDateChosen(chatID, messageID, barberID, date, serviceID, int(query.From.ID))
		}
	case "time":
		if len(parts) >= 5 {
			time := parts[4]
			h.handleBooking(chatID, messageID, query.ID, time, query.From)
		}
	case "booked":
		alert := tgbotapi.NewCallbackWithAlert(query.ID, "⚠️ Это время уже занято! Выберите другое.")
		h.bot.Request(alert)
		return
	}

}

func (h *Handler) HandleMessage(msg *tgbotapi.Message) {
	if msg == nil {
		return
	}
	clientID, err := h.getOrCreateClient(msg.From)
	if err != nil {
		log.Println(err)
	}
	_ = clientID

	message := tgbotapi.NewMessage(msg.Chat.ID, "Выберите интересующий вас раздел: ")
	message.ReplyMarkup = MainMenu()
	h.bot.Send(message)
	log.Printf("Пользователь %s написал: %s", msg.From.UserName, message.Text)
}

func (h *Handler) handleBarber(chatID int64, messageID int, barber string) {
	var text string
	var keyboards tgbotapi.InlineKeyboardMarkup

	switch barber {
	case "alex":
		text = texts.GetAlex()
		keyboards = AlexMenu()
	case "dmitrii":
		text = texts.GetDmitrii()
		keyboards = DmitriiMenu()
	}
	h.editMessage(chatID, messageID, text, keyboards)

}

func (h *Handler) handleBarberChosen(chatID int64, messageID int, barberID string, serviceID string, userID int) {
	state := h.getOrCreateState(userID)
	state.BarberID = barberID
	state.ServiceID = serviceID

	var serviceName string
	switch serviceID {
	case "haircut":
		serviceName = "Мужская стрижка"
	case "beard":
		serviceName = "Оформление бороды"
	case "combo":
		serviceName = "Комплекс"
	default:
		serviceName = serviceID
	}

	barberName := "Любой мастер"
	if barberID == "alex" {
		barberName = "Алекс (Top Barber)"
	} else if barberID == "dmitrii" {
		barberName = "Дмитрий"
	}

	text := fmt.Sprintf(
		"✂️ <b>Выбрано:</b> %s\n"+
			"👨‍🎨 <b>Мастер:</b> %s\n\n"+
			"📅 <i>Выберите удобный день для записи:</i>",
		serviceName, barberName,
	)
	h.editMessage(chatID, messageID, text, SelectDateMenu(serviceID, barberID))
}

func (h *Handler) handleService(chatID int64, messageID int, service string) {
	text := texts.HaircutService(service)
	h.editMessage(chatID, messageID, text, HaircutMenu(service))
}

func (h *Handler) editMessage(chatID int64, messageID int, text string, keyboards tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = &keyboards
	h.bot.Send(msg)
}
func (h *Handler) handleDateChosen(chatID int64, messageID int, barberID string, date string, serviceID string, userID int) {
	var serviceName string
	switch serviceID {
	case "haircut":
		serviceName = "Мужская стрижка"
	case "beard":
		serviceName = "Оформление бороды"
	case "combo":
		serviceName = "Комплекс"
	default:
		serviceName = serviceID
	}

	state := h.getOrCreateState(userID)
	state.Date = date

	barberName := "Любой мастер"
	if barberID == "alex" {
		barberName = "Алекс (Top Barber)"
	} else if barberID == "dmitrii" {
		barberName = "Дмитрий"
	}
	// Было: barberInt := mapBarberToID(barberName)
	// Стало:
	barberInt := mapBarberToID(barberID)
	bookedTimes, err := h.appointmentService.GetBookedTimes(barberInt, date)
	if err != nil {
		log.Printf("Error fetching booked times: %v", err)
		return
	}

	text := fmt.Sprintf(
		"✂️ <b>Услуга:</b> %s\n"+
			"👨‍🎨 <b>Мастер:</b> %s\n"+
			"📅 <b>Дата:</b> <b>%s</b>\n\n"+
			"🕒 <i>Выберите свободное время для записи:</i>",
		serviceName, barberName, date,
	)

	h.editMessage(chatID, messageID, text, SelectTimeMenu(serviceID, barberID, date, bookedTimes))
}

func (h *Handler) handleBooking(chatID int64, messageID int, id string, time string, tgUser *tgbotapi.User) {
	userID := int(tgUser.ID)

	state, exists := h.getState(userID)
	if !exists || state.BarberID == "" || state.Date == "" || state.ServiceID == "" {
		alert := tgbotapi.NewCallbackWithAlert(id, "⚠️ Сессия истекла. Начните запись заново.")
		h.bot.Request(alert)
		h.editMessage(chatID, messageID, texts.GetMainMenuText(), MainMenu())
		return
	}
	state.Time = time

	clientID, err := h.getOrCreateClient(tgUser)
	if err != nil {
		log.Printf("Error getting client for booking: %v", err)
		alert := tgbotapi.NewCallbackWithAlert(id, "⚠️ Ошибка при создании записи. Попробуйте позже.")
		h.bot.Request(alert)
		return
	}

	bookingDateTime := fmt.Sprintf("%s %s", state.Date, state.Time)
	masterInt := mapBarberToID(state.BarberID)
	ServiceInt := mapServiceToID(state.ServiceID)
	query := `INSERT INTO appointments (client_id, service_id, master_id, appointment_time, status) VALUES ($1, $2, $3, $4,'confirmed')`
	_, err = h.db.Exec(query, clientID, ServiceInt, masterInt, bookingDateTime)
	if err != nil {
		log.Printf("Error inserting appointment: %v", err)
		alert := tgbotapi.NewCallbackWithAlert(id, "⚠️ Не удалось сохранить запись в базу данных.")
		h.bot.Request(alert)
		return
	}

	h.resetState(int(tgUser.ID))

	alert := tgbotapi.NewCallbackWithAlert(id, "🎉 Запись успешно оформлена!")
	h.bot.Request(alert)
	h.editMessage(chatID, messageID, texts.GetMainMenuText(), MainMenu())
}

func (h *Handler) getOrCreateClient(tgUser *tgbotapi.User) (int, error) {
	if tgUser == nil {
		return 0, nil
	}
	return h.clientService.GetOrCreateClient(tgUser)
}

func (h *Handler) getOrCreateState(UserID int) *models.BookingState {
	h.mu.Lock()
	defer h.mu.Unlock()

	state, exists := h.userState[UserID]
	if !exists {
		state = &models.BookingState{}
		h.userState[UserID] = state
	}
	return state
}

func (h *Handler) getState(UserID int) (*models.BookingState, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	state, exists := h.userState[UserID]
	return state, exists
}

func (h *Handler) resetState(UserID int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.userState, UserID)
}

func mapBarberToID(barberID string) int {
	switch barberID {
	case "alex":
		return 1
	case "dmitrii":
		return 2
	default:
		return 1
	}
}

func mapServiceToID(serviceID string) int {
	switch serviceID {
	case "haircut":
		return 1
	case "beard":
		return 2
	case "combo":
		return 3
	default:
		return 1
	}
}
