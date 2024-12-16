package contacts

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/artem6554/task-9/internal/config"
	"github.com/artem6554/task-9/internal/db"
	"github.com/artem6554/task-9/internal/errors"
)

var dbConfig config.DbData

func init() {
	dbConfig = config.ReadDbConfig()
}

type Contact struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

func GetContact(name string) ([]byte, error) {
	var сontacts []Contact
	db, err := db.ConnectDB(dbConfig)
	if err != nil {
		return nil, errors.ErrInternalError
	}
	defer db.Close()
	queryString := fmt.Sprintf("select name, number from numbers where name = '%v'", name)
	rows, err := db.Query(queryString)
	if err != nil {
		return nil, errors.ErrInternalError
	}
	defer rows.Close()
	for rows.Next() {
		var contact Contact
		rows.Scan(&contact.Name, &contact.Number)
		сontacts = append(сontacts, contact)
	}
	result, err := json.Marshal(сontacts)
	if err != nil {
		return nil, errors.ErrInternalError
	}
	return result, nil
}

func AddContact(name string, number string) error {
	db, err := db.ConnectDB(dbConfig)
	if err != nil {
		return errors.ErrInternalError
	}
	defer db.Close()
	queryString := fmt.Sprintf("insert into numbers (name, number) values ('%v', '%v')", name, number)
	_, err = db.Query(queryString)
	if err != nil {
		return errors.ErrInternalError
	}
	return nil
}

func EditNumber(name string, number string) error {
	db, err := db.ConnectDB(dbConfig)
	if err != nil {
		return errors.ErrInternalError
	}
	defer db.Close()
	queryString := fmt.Sprintf("UPDATE numbers SET number = '%v' WHERE name = '%v'", number, name)
	_, err = db.Query(queryString)
	if err != nil {
		return errors.ErrInternalError
	}
	return nil
}

func DeleteContact(name string) error {
	db, err := db.ConnectDB(dbConfig)
	if err != nil {
		return errors.ErrInternalError
	}
	defer db.Close()
	queryString := fmt.Sprintf("DELETE FROM numbers WHERE name = '%v'", name)
	_, err = db.Query(queryString)
	if err != nil {
		return errors.ErrInternalError
	}
	return nil
}

func Exists(name string) error {
	db, err := db.ConnectDB(dbConfig)
	if err != nil {
		return errors.ErrInternalError
	}
	defer db.Close()
	queryString := fmt.Sprintf("select id, name, number from numbers where name = '%v'", name)
	rows, err := db.Query(queryString)
	if err != nil {
		return errors.ErrInternalError
	}
	defer rows.Close()
	if rows.Next() {
		return errors.ErrContactAlreadyExists
	}
	return nil
}

func CorrectNumber(number string) error {
	numberRegExp := regexp.MustCompile(`^(?:\+7|7|8)?\s*\(?(\d{3})\)?\s*(\d{3})\s*(\d{2})\s*(\d{2})$`)
	if numberRegExp.MatchString(number) {
		return nil
	}
	return errors.ErrIncorrectNumber
}
