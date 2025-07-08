package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

// -------------------- User CRUD --------------------

func (r *Repository) CreateUser(name string) (*User, error) {
	id := uuid.New()
	_, err := r.DB.Exec(`INSERT INTO users (id, name) VALUES ($1, $2)`, id, name)
	if err != nil {
		return nil, err
	}
	return &User{ID: id, Name: name}, nil
}

func (r *Repository) GetUser(id string) (*User, error) {
	var user User
	err := r.DB.QueryRow(`SELECT id, name FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(id string, name string) error {
	res, err := r.DB.Exec(`UPDATE users SET name = $1 WHERE id = $2`, name, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *Repository) DeleteUser(id string) error {
	res, err := r.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *Repository) GetUsers() ([]*User, error) {
	rows, err := r.DB.Query(`SELECT id, name FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// -------------------- Todo CRUD --------------------

func (r *Repository) CreateTodo(text string, done bool, userID string) (*Todo, error) {
	id := uuid.New()
	_, err := r.DB.Exec(`INSERT INTO todos (id, text, done, user_id) VALUES ($1, $2, $3, $4)`, id, text, done, userID)
	if err != nil {
		return nil, err
	}
	user, err := r.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return &Todo{ID: id, Text: text, Done: done, User: user}, nil
}

func (r *Repository) GetTodo(id string) (*Todo, error) {
	var todo Todo
	var userID string
	err := r.DB.QueryRow(`SELECT id, text, done, user_id FROM todos WHERE id = $1`, id).
		Scan(&todo.ID, &todo.Text, &todo.Done, &userID)
	if err != nil {
		return nil, err
	}
	user, err := r.GetUser(userID)
	if err != nil {
		return nil, err
	}
	todo.User = user
	return &todo, nil
}

func (r *Repository) UpdateTodo(id string, text string, done bool) error {
	res, err := r.DB.Exec(`UPDATE todos SET text = $1, done = $2 WHERE id = $3`, text, done, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("todo not found")
	}
	return nil
}

func (r *Repository) DeleteTodo(id string) error {
	res, err := r.DB.Exec(`DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("todo not found")
	}
	return nil
}
