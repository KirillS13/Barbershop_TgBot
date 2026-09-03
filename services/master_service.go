package services

import (
	"database/sql"
	"my-tg-bot/models"
)

type MasterService struct {
	db *sql.DB
}

func NewMasterService(db *sql.DB) *MasterService {
	return &MasterService{db: db}
}

func (m *MasterService) GetAllMaters() ([]models.Master, error) {
	var masters []models.Master
	rows, err := m.db.Query("SELECT * FROM masters WHERE is_active = true ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var master models.Master
		if err = rows.Scan(&master.Id, &master.Name, &master.Phone, &master.IsActive, &master.CreatedAt); err != nil {
			return nil, err
		}
		masters = append(masters, master)
	}
	return masters, nil
}

func (m *MasterService) GetOneMasterById(id int) (*models.Master, error) {
	var master models.Master
	err := m.db.QueryRow("SELECT * FROM master WHERE id = $1", id).Scan(&master.Id, &master.Name, &master.Phone, &master.IsActive, &master.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &master, nil
}
