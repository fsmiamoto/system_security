package main

import (
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/as/crypto"
	"github.com/fsmiamoto/system_security/kerberos/as/repository"
	"github.com/gofiber/fiber/v2"
)

const TGSKey = "5a6d29b8"

var TGSInitVector = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

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

		tgt, err := json.Marshal(
			contracts.TGT{
				ClientID:     client.ID,
				AccessPeriod: serviceReq.AccessPeriod,
				KeyClientTGS: key,
			})
		tgtBytes, _ := crypto.Encrypt([]byte(TGSKey), TGSInitVector, tgt)

		res, err :=
			json.Marshal(contracts.ASResponse{
				KeyClientTGS: key,
				Nonce:        serviceReq.Nonce,
			})

		resBytes, _ := crypto.Encrypt(client.SecretKey, client.InitVector, res)

		tgtRes := contracts.TGTResponse{
			CipheredASResponse: hex.EncodeToString(resBytes),
			CipheredTGT:        hex.EncodeToString(tgtBytes),
		}

		return c.JSON(tgtRes)
	})

	app.Post("/encrypt", func(c *fiber.Ctx) error {
		key := c.FormValue("key")
		iv := c.FormValue("iv")
		content := c.FormValue("content")

		data, err := crypto.Encrypt([]byte(key), []byte(iv), []byte(content))
		if err != nil {
			return err
		}

		return c.SendString(string(hex.EncodeToString(data)))
	})

	app.Listen(":3000")
}
