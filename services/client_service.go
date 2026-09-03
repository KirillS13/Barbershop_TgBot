package services

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ClientService struct {
	db *sql.DB
}

func NewClientService(db *sql.DB) *ClientService {
	return &ClientService{db: db}
}

func (s *ClientService) GetOrCreateClient(tgUser *tgbotapi.User) (int, error) {
	var clientID int
	query := `INSERT INTO clients (telegram_id, first_name ,username) VALUES ($1, $2, $3)
		ON CONFLICT (telegram_id) 
		DO UPDATE SET first_name = EXCLUDED.first_name, username = EXCLUDED.username
		RETURNING id;`
	err := s.db.QueryRow(query, tgUser.ID, tgUser.FirstName, tgUser.UserName).Scan(&clientID)
	if err != nil {
		return 0, err
	}
	return clientID, nil
}
