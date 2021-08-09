package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/as/repository"
	"github.com/fsmiamoto/system_security/kerberos/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var tgsKey, tgsInitVector string

func main() {
	app := fiber.New(fiber.Config{
		Prefork:   false,
		BodyLimit: 2 * 1024 * 1024,
	})

	app.Use(logger.New())

	app.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		request := contracts.TGTRequest{}

		if err := json.Unmarshal(c.Body(), &request); err != nil {
			log.Println(err)
			return err
		}

		client, err := repository.Get(request.ClientID)
		if err != nil {
			return err
		}

		data, err := hex.DecodeString(string(request.CipheredServiceRequest))
		if err != nil {
			return err
		}

		unencrypted, err := crypto.Decrypt(client.SecretKey, client.InitVector, data)

		serviceReq := new(contracts.ServiceRequest)
		if err := json.Unmarshal(unencrypted, serviceReq); err != nil {
			return err
		}

		key, err := crypto.GenKey()
		if err != nil {
			return err
		}

		tgt, err := json.Marshal(contracts.TGT{
			ClientID:     client.ID,
			AccessPeriod: serviceReq.AccessPeriod,
			CreatedAt:    time.Now(),
			KeyClientTGS: key,
		})

		tgtBytes, err := crypto.Encrypt([]byte(tgsKey), []byte(tgsInitVector), tgt)
		if err != nil {
			log.Fatal(err)
		}

		res, err := json.Marshal(contracts.ASResponse{
			KeyClientTGS:  key,
			TGSInitVector: tgsInitVector,
			Nonce:         serviceReq.Nonce,
		})

		resBytes, _ := crypto.Encrypt(client.SecretKey, client.InitVector, res)

		tgtRes := contracts.TGTResponse{
			CipheredASResponse: hex.EncodeToString(resBytes),
			CipheredTGT:        hex.EncodeToString(tgtBytes),
		}

		return c.JSON(tgtRes)
	})

	app.Listen(":3000")
}
