package main

import (
	"fmt"

	"github.com/0ne-zero/f4h/database"
	general_func "github.com/0ne-zero/f4h/utilities/functions/general"
	"github.com/0ne-zero/f4h/utilities/wrapper_logger"
	"github.com/0ne-zero/f4h/web/route"
)

func main() {
	if !general_func.IsUserRoot() {
		fmt.Println("Only root user can run this program (:\nProbably you forgot to use 'sudo' command.")
		//os.Exit(1)
	}

	db := database.InitializeOrGetDB()

	// If db is nil we kill the program, because we can't continue without database
	if db == nil {
		wrapper_logger.Fatal(&wrapper_logger.LogInfo{Message: "We really cannot connect to the database", ErrorLocation: general_func.GetCallerInfo(0)})
	}

	err := database.MigrateModels(db)

	print(err)
	//database.CreateTestData(db)
	//database.CreateEssentialData(db)
	// var user models.User
	// db.Debug().Preload("Votes").Preload("Comments").Preload("Activity").Preload("Addresses").Preload("Polls").Preload("Products").Preload("Carts").Preload("WalletInfos").Preload("WalletInfos").Preload("Orders").Preload("Roles").First(&user)
	// s, _ := models.IsExists(db, &user)
	// time level msg func file
	route := route.MakeRoute()
	route.Run(":8080")
	print("alive")
}
