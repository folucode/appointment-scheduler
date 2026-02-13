package db

import "log"

func (db *Database) FindUserByEmail(email string) (*User, error) {
	data, err := db.Read()

	if err != nil {
		return nil, err
	}

	for _, user := range data.Users {
		if user.Email == email {

			return &user, nil
		}
	}

	return nil, ErrUserNotFound
}

func (db *Database) CreateUser(user User) (User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	log.Printf("CreateUser....")

	data, err := db.read()

	if err != nil {
		return User{}, err
	}

	data.Users = append(data.Users, user)
	if err := db.save(data); err != nil {
		return User{}, err
	}

	return user, nil
}
