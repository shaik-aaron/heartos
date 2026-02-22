package intializers

import "github.com/shaik-aaron/fantasy-backend/models"

func MigrateDb() {
	DB.AutoMigrate(&models.User{}, &models.Session{})
}
