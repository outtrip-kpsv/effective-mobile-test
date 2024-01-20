package main

import (
  "em_test/internal/bl"
  "em_test/internal/config"
  "em_test/internal/db"
  "em_test/internal/io/http/handlers"
  "fmt"
  "github.com/gofiber/fiber/v2"
  "log"
)

func main() {
  conf, err := config.InitConfServ()
  if err != nil {
    panic("cannot initialize config")
  }
  fmt.Println(conf)
  app := fiber.New()
  dbPg, _ := db.NewDB(conf.Logger, conf.Options)
  dbRepo := db.NewDBRepo(dbPg)

  blInst := bl.NewBL(dbRepo, conf.Logger)
  x := handlers.NewController(blInst, conf.Logger)

  app.Post("/create", x.CreatePerson)
  app.Delete("/del/:id", x.DeletePerson)
  app.Get("/people", x.GetPeople)
  app.Patch("/update/:id", x.UpdatePerson)

  log.Fatal(app.Listen(":3000"))

  //
  //// REST методы
  //app.Post("/people", createPerson)
  //app.Get("/people", getPeople)
  //app.Get("/people/:id", getPerson)
  //app.Put("/people/:id", updatePerson)
  //app.Delete("/people/:id", deletePerson)
}
