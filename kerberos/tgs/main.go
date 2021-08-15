package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	as "github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/crypto"
	"github.com/fsmiamoto/system_security/kerberos/tgs/contracts"
	"github.com/fsmiamoto/system_security/kerberos/tgs/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	secretKey, initVector string
	accessPeriod          = 10 * time.Minute
)

var (
	ErrInvalidRequest  = fiber.NewError(http.StatusBadRequest, "invalid service ticket request")
	ErrServiceNotFound = fiber.NewError(http.StatusNotFound, "service not found")
	ErrInternalError   = fiber.NewError(http.StatusInternalServerError, "internal server error")
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:   false,
		BodyLimit: 2 * 1024 * 1024,
	})

	app.Use(logger.New())

	app.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		svcTicketReq := &contracts.ServiceTicketRequest{}
		if err := json.Unmarshal(c.Body(), svcTicketReq); err != nil {
			log.Debug().Err(err).Msg("error while unmarshaling ticker request")
			return ErrInvalidRequest
		}

		tgtBytes, err := hex.DecodeString(svcTicketReq.CipheredTGT)
		if err != nil {
			log.Debug().Err(err).Msg("error while decoding hex of tgt")
			return fiber.NewError(fiber.StatusBadRequest, "invalid service ticket request")
		}

		tgtBytes, err = crypto.Decrypt([]byte(secretKey), []byte(initVector), tgtBytes)
		if err != nil {
			log.Debug().Err(err).Msg("error while decrypting tgt")
			return fiber.NewError(fiber.StatusBadRequest, "invalid service ticket request")
		}

		tgt := &as.TGT{}

		if err := json.Unmarshal(tgtBytes, tgt); err != nil {
			log.Debug().Err(err).Msg("error while unmarshaling tgt")
			return ErrInvalidRequest
		}

		if time.Now().After(tgt.CreatedAt.Add(tgt.AccessPeriod)) {
			log.Info().Msgf("expired tgt from %v", tgt.ClientID)
			return ErrInvalidRequest
		}

		svcReqBytes, err := hex.DecodeString(svcTicketReq.CipheredServiceRequest)
		if err != nil {
			log.Debug().Err(err).Msg("error while decoding hex from ciphered service request")
			return ErrInvalidRequest
		}

		svcReqBytes, err = crypto.Decrypt([]byte(tgt.KeyClientTGS), []byte(initVector), svcReqBytes)
		if err != nil {
			log.Debug().Err(err).Msgf("error while decrypting request")
			return ErrInvalidRequest
		}

		svcReq := &contracts.ServiceRequest{}
		if err := json.Unmarshal(svcReqBytes, svcReq); err != nil {
			log.Debug().Err(err).Msgf("error while unmarshaling request")
			return ErrInvalidRequest
		}

		svc, err := repository.Get(svcReq.ServiceID)
		if err != nil {
			log.Debug().Err(err).Msg("service not found")
			return ErrServiceNotFound
		}

		serviceClientKey, _ := crypto.GenKey()

		serviceTicket := &contracts.ServiceTicket{
			ClientID:         svcReq.ClientID,
			AccessPeriod:     accessPeriod,
			KeyClientService: serviceClientKey,
		}

		serviceTicketBytes, err := json.Marshal(serviceTicket)
		if err != nil {
			log.Error().Err(err).Msgf("error while marshaling ticket")
			return ErrInternalError
		}

		serviceTicketBytes, err = crypto.Encrypt(svc.SecretKey, svc.InitVector, serviceTicketBytes)
		if err != nil {
			log.Error().Err(err).Msgf("error while encrypting service ticket")
			return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}

		tgsRes := &contracts.TGSResponse{
			AccessPeriod:      accessPeriod,
			KeyClientService:  serviceClientKey,
			Nonce:             svcReq.Nonce,
			ServiceInitVector: string(svc.InitVector),
		}

		tgsResBytes, err := json.Marshal(tgsRes)
		if err != nil {
			log.Debug().Err(err).Msg("error while marshaling response")
			return ErrInternalError
		}

		tgsResBytes, err = crypto.Encrypt([]byte(tgt.KeyClientTGS), []byte(initVector), tgsResBytes)
		if err != nil {
			log.Error().Err(err).Msgf("error while encrypting response")
			return ErrInternalError
		}

		svcTicketRes := &contracts.ServiceTicketResponse{
			CipheredServiceTicket: hex.EncodeToString(serviceTicketBytes),
			CipheredTGSResponse:   hex.EncodeToString(tgsResBytes),
		}

		return c.JSON(svcTicketRes)
	})

	app.Listen(":4000")
}
