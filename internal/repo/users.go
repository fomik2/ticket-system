package repo

import (
	"database/sql"
	"fmt"

	"github.com/fomik2/ticket-system/internal/entities"
)

func (t *Repo) GetUser(id int) (entities.Users, error) {
	row := t.db.QueryRow("SELECT * FROM users WHERE id=?;", id)
	user := entities.Users{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return entities.Users{}, err
	}
	return user, nil
}

func (t *Repo) ListUsers() ([]entities.Users, error) {
	return []entities.Users{}, nil
}

func (t *Repo) CreateUser(user entities.Users) (entities.Users, error) {
	checkIfExist, err := t.FindUser(user.Name)
	if err != nil {
		return entities.Users{}, err
	}
	if checkIfExist.Name != "" {
		return entities.Users{}, fmt.Errorf("user already exist")
	}
	_, err = t.db.Exec("INSERT INTO users VALUES(NULL,?,?,?,?);", user.Name, user.Password, user.Email, user.CreatedAt.Format("2006-01-02 15:04"))
	if err != nil {
		return entities.Users{}, err
	}
	return user, nil
}

func (t *Repo) DeleteUser(id int) error {
	_, err := t.db.Exec("DELETE FROM users WHERE id=?;", id)
	if err != nil {
		return err
	}
	return nil
}

func (t *Repo) FindUser(username string) (entities.Users, error) {
	row := t.db.QueryRow("SELECT * FROM users WHERE name=?;", username)
	user := entities.Users{}
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return entities.Users{}, nil
	}
	if err != nil {
		return entities.Users{}, err
	}
	return user, nil
}

func (t *Repo) UpdateUser(user entities.Users) (entities.Users, error) {
	_, err := t.db.Exec("UPDATE users SET name=?, email=?, password=? WHERE id=?;", user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		return entities.Users{}, err
	}
	return entities.Users{}, nil
}

func (t *Repo) IsUserExistInDB(username string) (int, error) {
	row := t.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name=?);", username)
	var isUserExist int
	err := row.Scan(&isUserExist)
	if err != nil {
		return isUserExist, err
	}
	return isUserExist, nil
}
