package services

import "database/sql"

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
