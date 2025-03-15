package main

import (
	p "github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/JneiraS/AMS/src/infrastructure/web/views"
)

func main() {
	db := p.CreateDB()

	// Lancement du worker pour traiter les transactions en arrière-plan
	go p.Worker(db)

	views.SetupRouter(db)

}
