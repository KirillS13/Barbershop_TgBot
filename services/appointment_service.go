package services

import (
	"context"
	"database/sql"
	"my-tg-bot/models"
)

type AppointmentService struct {
	db *sql.DB
}

func NewAppointmentService(db *sql.DB) *AppointmentService {
	return &AppointmentService{db: db}
}

func (appointmentService *AppointmentService) GetBookedTimes(barberID int, date string) ([]string, error) {
	var bookedTimes []string

	query := `SELECT TO_CHAR(appointment_time, 'HH24:MI') FROM appointments WHERE master_id=$1 AND DATE(appointment_time)=$2`
	rows, err := appointmentService.db.Query(query, barberID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bookedTime string
		if err = rows.Scan(&bookedTime); err != nil {
			return nil, err
		}
		bookedTimes = append(bookedTimes, bookedTime)
	}
	return bookedTimes, nil
}

func (appointmentService *AppointmentService) GetAppointments(clientId int64) ([]models.Appointment, error) {
	var appointments []models.Appointment
	var id int64

	req := `SELECT id FROM clients WHERE telegram_id=$1`
	err := appointmentService.db.QueryRowContext(context.Background(), req, clientId).Scan(&id)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, client_id, service_id, master_id, appointment_time, status, created_at FROM appointments WHERE client_id=$1`
	rows, err := appointmentService.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var appointment models.Appointment
		if err = rows.Scan(&appointment.Id,
			&appointment.ClientId,
			&appointment.ServerId, // сюда запишется service_id
			&appointment.MasterId,
			&appointment.Time,
			&appointment.Status,
			&appointment.CreatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}
