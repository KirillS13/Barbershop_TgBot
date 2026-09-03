package main

import (
	"log"
	"my-tg-bot/services"
	"my-tg-bot/services/database"
	"my-tg-bot/telegram"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	clientService := services.NewClientService(db)
	appointmentService := services.NewAppointmentService(db)

	updates := bot.GetUpdatesChan(u)
	h := telegram.NewHandler(bot, db, clientService, appointmentService)

	for update := range updates {
		if update.CallbackQuery != nil {
			h.HandleCallBack(update.CallbackQuery)
			continue
		}
		if update.Message != nil {
			h.HandleMessage(update.Message)
		}
	}

}
