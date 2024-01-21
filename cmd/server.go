package main

import (
	"em_test/internal/bl"
	"em_test/internal/config"
	"em_test/internal/db"
	"em_test/internal/io"
	"em_test/internal/io/http/handlers"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf, err := config.InitConfServ()
	if err != nil {
		panic("cannot initialize config")
	}

	app := fiber.New()
	dbPg, _ := db.NewDB(conf.Logger, conf.Options)
	dbRepo := db.NewDBRepo(dbPg)

	blInst := bl.NewBL(dbRepo, conf.Logger)
	controller := handlers.NewController(blInst, conf.Logger)
	io.SetupRoutes(app, *controller)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Fatal(app.Listen(fmt.Sprintf("%s:%s", conf.Options.Host, conf.Options.Port)))
	}()

	<-stop

	if err := app.Shutdown(); err != nil {
		fmt.Println("Error while shutting down:", err)
	}
	fmt.Println("Server gracefully stopped")
}
