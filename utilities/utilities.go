package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/bcrypt"
)

func ReadSettingFile(setting_path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	file, err := ioutil.ReadFile(setting_path)
	if err != nil {
		return nil, errors.New("Error when opening setting file")
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, errors.New("Error during unmarshal setting file")
	}
	return data, nil
}
func ReadFieldInSettingFile(setting_path string, field_name string) (string, error) {
	settings, err := ReadSettingFile(setting_path)
	if err != nil {
		return "", err
	}

	field_value, exists := settings[field_name]
	if !exists {
		return "", errors.New("Field not exists in setting file")
	}

	// Convert type interface to string and return value of field
	return fmt.Sprint(field_value), nil

}
func HashPassword(pass string) (string, error) {
	// Generate bcrypt hash from password with 17 cost
	hash_bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 17)
	if err != nil {
		return "", err
	}
	return string(hash_bytes), nil
}
func ComparePassword(hashed_pass string, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed_pass), []byte(pass))
}
