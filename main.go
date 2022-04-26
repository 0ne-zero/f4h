package main

import (
	"github.com/0ne-zero/f4h/database"
	"github.com/0ne-zero/f4h/web/route"
)

func main() {
	db, err := database.Initialize()
	if err != nil {

	}

	err = database.MigrateModels(db)
	// database.CreateTestData(db)

	// database.CreateEssentialData(db)
	// var user models.User
	// db.Debug().Preload("Votes").Preload("Comments").Preload("Activity").Preload("Addresses").Preload("Polls").Preload("Products").Preload("Carts").Preload("WalletInfos").Preload("WalletInfos").Preload("Orders").Preload("Roles").First(&user)
	// s, _ := models.IsExists(db, &user)
	// time level msg func file

	route := route.MakeRoute()
	route.Run(":8080")
}
