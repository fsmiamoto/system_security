package main

import (
	"encoding/json"
	"log"

	"github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/as/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:   false,
		BodyLimit: 2 * 1024 * 1024,
	})

	app.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		request := contracts.TGTRequest{}

		if err := json.Unmarshal(c.Body(), &request); err != nil {
			log.Println(err)
			return err
		}

		client, err := repository.Get(request.ClientID)
		if err != nil {
			log.Println(err)
			return err
		}

		return c.SendString(client.ID)
	})

	app.Listen(":3000")
}
