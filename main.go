package main

import (
	"fmt"

	"github.com/0ne-zero/f4h/database"
	"github.com/0ne-zero/f4h/web/route"
)

func main() {
	db, err := database.Initialize()
	if err != nil {

	}
	err = database.MigrateModels(db)
	//database.CreateTestData(db)
	// database.CreateEssentialData(db)
	// var user models.User
	// db.Debug().Preload("Votes").Preload("Comments").Preload("Activity").Preload("Addresses").Preload("Polls").Preload("Products").Preload("Carts").Preload("WalletInfos").Preload("WalletInfos").Preload("Orders").Preload("Roles").First(&user)
	// s, _ := models.IsExists(db, &user)

	route := route.MakeRoute()
	route.Run(":8080")
	fmt.Println("hello world")

}
