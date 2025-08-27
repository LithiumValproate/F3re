package user

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
)

type UserType string

const (
	TypeStudent UserType = "student"
	TypeTeacher UserType = "teacher"
	TypeAdmin   UserType = "admin"
	TypeBot     UserType = "bot"
)

type User interface {
	ID() string
	Name() string
	PasswordHash() string
	CheckPassword(password string) bool
	Type() UserType
}

type baseUser struct {
	id           string
	name         string
	passwordHash string
}

func newBaseUser(id, name, password string) (*baseUser, error) {
	u := &baseUser{
		id:   id,
		name: name,
	}
	err := u.SetPassword(password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *baseUser) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.passwordHash = string(hashedPassword)
	return nil
}

func (u *baseUser) ID() string {
	return u.id
}

func (u *baseUser) Name() string {
	return u.name
}

func (u *baseUser) PasswordHash() string {
	return u.passwordHash
}

func (u *baseUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
	return err == nil
}

func (u *baseUser) marshalJSON(userType UserType) ([]byte, error) {
	return json.Marshal(&struct {
		ID       string   `json:"id"`
		Name     string   `json:"name"`
		UserType UserType `json:"user_type"`
	}{
		ID:       u.ID(),
		Name:     u.Name(),
		UserType: userType,
	})
}

type Student struct {
	*baseUser
}

func NewStudent(id, name, password string) (*Student, error) {
	u, err := newBaseUser(id, name, password)
	if err != nil {
		return nil, err
	}
	return &Student{u}, nil
}

func NewStudentFromDB(id, name, passwordHash string) (*Student, error) {
	u := &baseUser{
		id:           id,
		name:         name,
		passwordHash: passwordHash,
	}
	return &Student{u}, nil
}

func (s *Student) Type() UserType {
	return TypeStudent
}

func (s *Student) MarshalJSON() ([]byte, error) {
	return s.marshalJSON(s.Type())
}

type Teacher struct {
	*baseUser
}

func NewTeacher(id, name, password string) (*Teacher, error) {
	u, err := newBaseUser(id, name, password)
	if err != nil {
		return nil, err
	}
	return &Teacher{u}, nil
}

func NewTeacherFromDB(id, name, passwordHash string) (*Teacher, error) {
	u := &baseUser{
		id:           id,
		name:         name,
		passwordHash: passwordHash,
	}
	return &Teacher{u}, nil
}

func (t *Teacher) Type() UserType {
	return TypeTeacher
}

func (t *Teacher) MarshalJSON() ([]byte, error) {
	return t.marshalJSON(t.Type())
}

type Admin struct {
	*baseUser
}

func NewAdmin(id, name, password string) (*Admin, error) {
	u, err := newBaseUser(id, name, password)
	if err != nil {
		return nil, err
	}
	return &Admin{u}, nil
}

func NewAdminFromDB(id, name, passwordHash string) (*Admin, error) {
	u := &baseUser{
		id:           id,
		name:         name,
		passwordHash: passwordHash,
	}
	return &Admin{u}, nil
}

func (a *Admin) Type() UserType {
	return TypeAdmin
}

func (a *Admin) MarshalJSON() ([]byte, error) {
	return a.marshalJSON(a.Type())
}

type BotUser struct {
	*baseUser
}

func NewBotUser(id, name, password string) (*BotUser, error) {
	u, err := newBaseUser(id, name, password)
	if err != nil {
		return nil, err
	}
	return &BotUser{u}, nil
}

func NewBotUserFromDB(id, name, passwordHash string) (*BotUser, error) {
	u := &baseUser{
		id:           id,
		name:         name,
		passwordHash: passwordHash,
	}
	return &BotUser{u}, nil
}

func (b *BotUser) Type() UserType {
	return TypeBot
}

func (b *BotUser) MarshalJSON() ([]byte, error) {
	return b.marshalJSON(b.Type())
}
