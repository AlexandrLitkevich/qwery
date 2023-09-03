package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/user"
	"github.com/AlexandrLitkevich/qwery/internal/storage.go"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

//TODO: Edit ETCD

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New" //Operation this name function

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//Pet migration
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS user(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		position TEXT NOT NULL);
	

`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

//TODO: DeleteUrl

func (s *Storage) CreateUser(userInfo user.Request) (*user.User, error) {
	const op = "storage.sqlite.CreateUser"

	//Generate ID
	//TODO change type int on string
	userId, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	us := userId.String()
	fmt.Println("this US ===", us)
	// TODO: asyn

	stmt, err := s.db.Prepare("INSERT INTO user(id, name, age, position) VALUES(?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(us, userInfo.Name, userInfo.Age, userInfo.Position)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createdUser, err := s.GetUser(us)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get last createdUser id: %w", op, err)
	}

	return createdUser, nil
}

func (s *Storage) GetUser(userId string) (*user.User, error) {
	const op = "storage.sqlite.GetUser"
	stmt, err := s.db.Prepare("SELECT * FROM user WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var (
		ID       string
		Name     string
		Age      int
		Position string
	)

	err = stmt.QueryRow(userId).Scan(&ID, &Name, &Age, &Position)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return &user.User{ID: ID, Name: Name, Age: Age, Position: Position}, nil
}
