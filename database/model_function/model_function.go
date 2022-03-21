package model_function

import (
	"fmt"
	"time"

	"github.com/0ne-zero/f4h/database"
	"github.com/0ne-zero/f4h/database/model"
	"github.com/0ne-zero/f4h/utilities/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var db *gorm.DB

type Model interface {
	model.User | model.Product_Category | model.Product | model.Request
}

func init() {
	var err error
	db, err = database.Initialize()
	if err != nil {
		log.Log(logrus.Fatal, err)
	}
}

func Add[m Model](model *m) error {
	return db.Create(model).Error
}
func Get[m Model](model *[]m, limit int, orderBy string, orderType string) error {
	var err error
	if limit < 1 {
		err = db.Order(fmt.Sprintf("%s %s", orderBy, orderType)).Find(model).Error
	} else {
		err = db.Order(fmt.Sprintf("%s %s", orderBy, orderType)).Find(model).Limit(limit).Error
	}
	return err
}
func IsExistsByID[m Model](model *m, id uint) (bool, error) {
	var exists bool
	err := db.Model(model).Select("count(*) >0").Where("ID = ?", id).Find(&exists).Error
	return exists, err
}
func GetByID[m Model](model *m, id uint) error {
	return db.Where("id = ?", id).First(model).Error
}
func IsUserExistsByUsername(username string) (bool, error) {
	var exists bool
	err := db.Model(&model.User{}).Select("count(*) >0").Where("username = ?", username).Find(&exists).Error
	return exists, err
}
func GetUserByUsername(username string) (model.User, error) {
	var u model.User
	err := db.Where("username = ?", username).First(&u).Error
	return u, err

}
func GetUserIDByUsername(username string) (uint, error) {
	var u model.User
	err := db.Select("id").Where("username = ?", username).First(&u).Error
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

func TooManyRequest(ip string, url string, method string) (bool, error) {
	var req_count int64
	now := time.Now().UTC()
	one_hour_ago := now.Add(time.Duration(-1) * time.Hour)
	err := db.Model(&model.Request{}).Where("ip = ? AND url = ? AND method = ? AND time <= ? ", ip, url, method, one_hour_ago).Count(&req_count).Error
	if err != nil {
		return false, err
	}
	if req_count > 100 {
		return true, nil
	}
	return false, nil
}
