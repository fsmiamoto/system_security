package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	as "github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/crypto"
	"github.com/fsmiamoto/system_security/kerberos/tgs/contracts"
	"github.com/fsmiamoto/system_security/kerberos/tgs/repository"
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

		svcTicketReq := &contracts.ServiceTicketRequest{}
		mustBeNil(json.Unmarshal(c.Body(), svcTicketReq))

		tgtBytes, err := hex.DecodeString(svcTicketReq.CipheredTGT)
		mustBeNil(err)

		tgtBytes, err = crypto.Decrypt([]byte(secretKey), []byte(initVector), tgtBytes)

		tgt := &as.TGT{}

		mustBeNil(json.Unmarshal(tgtBytes, tgt))

		fmt.Println("TGT:")
		fmt.Println("Client ID:", tgt.ClientID)
		fmt.Println("Key:", tgt.KeyClientTGS)
		fmt.Println("Period:", tgt.AccessPeriod)

		svcReqBytes, err := hex.DecodeString(svcTicketReq.CipheredServiceRequest)
		mustBeNil(err)

		svcReqBytes, err = crypto.Decrypt([]byte(tgt.KeyClientTGS), []byte(initVector), svcReqBytes)

		svcReq := &contracts.ServiceRequest{}
		mustBeNil(json.Unmarshal(svcReqBytes, svcReq))

		svc, err := repository.Get(svcReq.ServiceID)
		if err != nil {
			fmt.Println(err)
			return fiber.NewError(fiber.StatusBadRequest, "invalid service")
		}

		serviceClientKey, err := crypto.GenKey()

		serviceTicket := &contracts.ServiceTicket{
			ClientID:         svcReq.ClientID,
			AccessPeriod:     accessPeriod,
			KeyClientService: serviceClientKey,
		}

		serviceTicketBytes, err := json.Marshal(serviceTicket)
		mustBeNil(err)

		serviceTicketBytes, err = crypto.Encrypt(svc.SecretKey, svc.InitVector, serviceTicketBytes)
		mustBeNil(err)

		tgsRes := &contracts.TGSResponse{
			AccessPeriod:     accessPeriod,
			KeyClientService: serviceClientKey,
			Nonce:            svcReq.Nonce,
		}

		tgsResBytes, err := json.Marshal(tgsRes)
		mustBeNil(err)
		tgsResBytes, err = crypto.Encrypt([]byte(tgt.KeyClientTGS), []byte(initVector), tgsResBytes)
		mustBeNil(err)

		svcTicketRes := &contracts.ServiceTicketResponse{
			CipheredServiceTicket: hex.EncodeToString(serviceTicketBytes),
			CipheredTGSResponse:   hex.EncodeToString(tgsResBytes),
		}

		return c.JSON(svcTicketRes)
	})

	app.Listen(":4000")
}

func mustBeNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
