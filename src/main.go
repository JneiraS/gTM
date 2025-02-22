package main

import (
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/JneiraS/AMS/src/infrastructure/web/views"
)

func main() {
	db := persistence.CreateDB()
	views.SetupRouter(db)

}
