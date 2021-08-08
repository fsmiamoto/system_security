package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/fsmiamoto/system_security/kerberos/crypto"
	"github.com/fsmiamoto/system_security/kerberos/service/contracts"
	tgs "github.com/fsmiamoto/system_security/kerberos/tgs/contracts"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var secretKey, initVector string

var accessPeriod = 10 * time.Minute

func main() {
	app := fiber.New(fiber.Config{
		Prefork:   false,
		BodyLimit: 2 * 1024 * 1024,
	})

	app.Use(logger.New())

	app.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		req := &contracts.ServiceRequest{}
		mustBeNil(json.Unmarshal(c.Body(), req))

		ticketBytes, err := hex.DecodeString(req.CipheredServiceTicket)
		mustBeNil(err)
		ticketBytes, err = crypto.Decrypt([]byte(secretKey), []byte(initVector), ticketBytes)

		t := &tgs.ServiceTicket{}
		mustBeNil(json.Unmarshal(ticketBytes, t))

		res := &contracts.Response{
			Result: time.Now(),
			Nonce:  rand.Uint64(),
		}

		resBytes, err := json.Marshal(res)
		mustBeNil(err)

		resBytes, err = crypto.Encrypt([]byte(t.KeyClientService), []byte(initVector), resBytes)

		resp := &contracts.ServiceResponse{
			CipheredResponse: hex.EncodeToString(resBytes),
		}
		return c.JSON(resp)
	})

	app.Listen(":5000")
}

func mustBeNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
