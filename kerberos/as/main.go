package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/as/repository"
	"github.com/fsmiamoto/system_security/kerberos/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/rs/zerolog/log"
)

var tgsKey, tgsInitVector string

var (
	ErrInvalidRequest = fiber.NewError(http.StatusBadRequest, "invalid request")
	ErrClientNotFound = fiber.NewError(http.StatusNotFound, "client not found")
	ErrInternalError  = fiber.NewError(http.StatusInternalServerError, "internal server error")
)

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
			log.Debug().Err(err).Msgf("error while unmarshaling request")
			return ErrInvalidRequest
		}

		client, err := repository.Get(request.ClientID)
		if err != nil {
			log.Debug().Err(err).Msgf("client not found")
			return ErrClientNotFound
		}

		data, err := hex.DecodeString(string(request.CipheredServiceRequest))
		if err != nil {
			log.Debug().Err(err).Msgf("error while hex decoding")
			return ErrInvalidRequest
		}

		unencrypted, err := crypto.Decrypt(client.SecretKey, client.InitVector, data)
		if err != nil {
			log.Debug().Err(err).Msgf("error while decrypting request")
			return ErrInvalidRequest
		}

		serviceReq := new(contracts.ServiceRequest)
		if err := json.Unmarshal(unencrypted, serviceReq); err != nil {
			log.Debug().Err(err).Msgf("error while unmarshaling request")
			return ErrInvalidRequest
		}

		key, _ := crypto.GenKey()

		tgt, err := json.Marshal(contracts.TGT{
			ClientID:     client.ID,
			AccessPeriod: serviceReq.AccessPeriod,
			CreatedAt:    time.Now(),
			KeyClientTGS: key,
		})
		if err != nil {
			log.Error().Err(err).Msgf("error while marshaling tgt")
			return ErrInternalError
		}

		tgtBytes, err := crypto.Encrypt([]byte(tgsKey), []byte(tgsInitVector), tgt)
		if err != nil {
			log.Error().Err(err).Msgf("error while encrypting tgt")
			return ErrInternalError
		}

		res, err := json.Marshal(contracts.ASResponse{
			KeyClientTGS:  key,
			TGSInitVector: tgsInitVector,
			Nonce:         serviceReq.Nonce,
		})
		if err != nil {
			log.Error().Err(err).Msgf("error while marshaling response")
			return ErrInternalError
		}

		resBytes, err := crypto.Encrypt(client.SecretKey, client.InitVector, res)
		if err != nil {
			log.Error().Err(err).Msgf("error while encrypting response")
			return ErrInternalError
		}

		tgtRes := contracts.TGTResponse{
			CipheredASResponse: hex.EncodeToString(resBytes),
			CipheredTGT:        hex.EncodeToString(tgtBytes),
		}

		return c.JSON(tgtRes)
	})

	app.Listen(":3000")
}
