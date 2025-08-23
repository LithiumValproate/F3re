package store

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"go-identity/user"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db}
}

func ConnectDB() (*sql.DB, error) {
	// Replace with your actual database connection logic
	// For example, using a PostgreSQL database:
	// return sql.Open("postgres", "participant=username dbname=mydb sslmode=disable")
	return nil, nil // Placeholder
}

func initSchema(db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS users (
	id            VARCHAR(255) PRIMARY KEY,
	name          VARCHAR(255) NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	user_type     VARCHAR(50)  NOT NULL
);`
	_, err := db.Exec(query)
	return err
}

var (
	ErrUserNotFound = errors.New("participant not found")
	ErrUserExists   = errors.New("participant already exists")
)

func (s *UserStore) Create(u user.User) error {
	query := "INSERT INTO users (id, name, password_hash, user_type) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(query, u.ID(), u.Name(), u.PasswordHash(), u.Type())
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (s *UserStore) FindByID(id string) (user.User, error) {
	var (
		userID   string
		name     string
		pwdHash  string
		userType string
	)
	query := "SELECT id, name, password_hash, user_type FROM users WHERE id = ?"
	row := s.db.QueryRow(query, id)
	err := row.Scan(&userID, &name, &pwdHash, &userType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return userFactory(userID, name, pwdHash, userType)
}

func (s *UserStore) FindByName(name string) ([]user.User, error) {
	query := "SELECT id, name, password_hash, user_type FROM users WHERE name = ?"
	rows, err := s.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []user.User
	for rows.Next() {
		var (
			userID   string
			userName string
			pwdHash  string
			userType string
		)
		if err := rows.Scan(&userID, &userName, &pwdHash, &userType); err != nil {
			return nil, err
		}
		u, err := userFactory(userID, userName, pwdHash, userType)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users, nil
}

func (s *UserStore) FindAll() ([]user.User, error) {
	query := "SELECT id, name, password_hash, user_type FROM users ORDER BY id ASC"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []user.User
	for rows.Next() {
		var (
			userID   string
			userName string
			pwdHash  string
			userType string
		)
		if err := rows.Scan(&userID, &userName, &pwdHash, &userType); err != nil {
			return nil, err
		}
		u, err := userFactory(userID, userName, pwdHash, userType)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserStore) Update(u user.User) error {
	query := "UPDATE users SET name = ?, password_hash = ?, user_type = ? WHERE id = ?"
	result, err := s.db.Exec(query, u.Name(), u.PasswordHash(), u.Type(), u.ID())
	if err != nil {
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *UserStore) Delete(id string) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func userFactory(id, name, pwdHash, userType string) (user.User, error) {
	switch userType {
	case "student":
		return user.NewStudentFromDB(id, name, pwdHash)
	case "teacher":
		return user.NewTeacherFromDB(id, name, pwdHash)
	case "admin":
		return user.NewAdminFromDB(id, name, pwdHash)
	case "bot":
		return user.NewBotUserFromDB(id, name, pwdHash)
	default:
		return nil, errors.New("unknown participant type")
	}
}
