package main

import (
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/fsmiamoto/system_security/kerberos/crypto"
	"github.com/fsmiamoto/system_security/kerberos/service/contracts"
	tgs "github.com/fsmiamoto/system_security/kerberos/tgs/contracts"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
)

var (
	secretKey, initVector string
	accessPeriod          = 10 * time.Minute
)

var (
	ErrInvalidRequest = fiber.NewError(http.StatusBadRequest, "invalid request")
	ErrUnauthorized   = fiber.NewError(http.StatusUnauthorized, "unauthorized")
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

		req := &contracts.ServiceRequest{}
		if err := json.Unmarshal(c.Body(), req); err != nil {
			log.Debug().Err(err).Msgf("invalid request")
			return ErrInvalidRequest
		}

		ticketBytes, err := hex.DecodeString(req.CipheredServiceTicket)
		if err := json.Unmarshal(c.Body(), req); err != nil {
			log.Debug().Err(err).Msgf("invalid request")
			return ErrInvalidRequest
		}
		ticketBytes, err = crypto.Decrypt([]byte(secretKey), []byte(initVector), ticketBytes)
		if err := json.Unmarshal(c.Body(), req); err != nil {
			log.Debug().Err(err).Msgf("invalid request")
			return ErrInvalidRequest
		}

		t := &tgs.ServiceTicket{}
		if err := json.Unmarshal(ticketBytes, t); err != nil {
			log.Debug().Err(err).Msgf("invalid request")
			return ErrInvalidRequest
		}

		res := &contracts.Response{
			Result: time.Now(),
			Nonce:  rand.Uint64(),
		}

		if t.CreatedAt.Add(t.AccessPeriod).After(time.Now()) {
			log.Info().Msgf("ticket expired, created at: %v", t.CreatedAt)
			return ErrUnauthorized
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			log.Error().Err(err).Msgf("error while trying to marshal json")
			return ErrInternalError
		}

		resBytes, err = crypto.Encrypt([]byte(t.KeyClientService), []byte(initVector), resBytes)
		if err != nil {
			log.Error().Err(err).Msgf("error while trying to marshal json")
			return ErrInternalError
		}

		resp := &contracts.ServiceResponse{
			CipheredResponse: hex.EncodeToString(resBytes),
		}

		return c.JSON(resp)
	})

	app.Listen(":5000")
}
