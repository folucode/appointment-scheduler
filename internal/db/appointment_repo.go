package db

import (
	"errors"
)

func (db *Database) CreateAppointment(appt Appointment) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := db.read()
	if err != nil {
		return err
	}

	data.Appointments = append(data.Appointments, appt)

	return db.save(data)
}

func (db *Database) GetAppointment(id string) (Appointment, error) {
	data, err := db.Read()
	if err != nil {
		return Appointment{}, err
	}

	for _, a := range data.Appointments {
		if a.ID == id {
			return a, nil
		}
	}

	return Appointment{}, errors.New("appointment not found")
}

func (db *Database) GetAppointments(user_id string) ([]Appointment, error) {
	data, err := db.Read()
	if err != nil {
		return []Appointment{}, err
	}

	result := []Appointment{}

	for _, a := range data.Appointments {
		if a.UserID == user_id {
			result = append(result, a)
		}
	}

	return result, nil
}

func (db *Database) DeleteAppointment(id string) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := db.read()
	if err != nil {
		return false, err
	}

	for i, a := range data.Appointments {
		if a.ID == id {
			data.Appointments = append(
				data.Appointments[:i],
				data.Appointments[i+1:]...,
			)

			if err := db.Save(data); err != nil {
				return false, err
			}

			return true, nil
		}
	}

	return false, nil
}
