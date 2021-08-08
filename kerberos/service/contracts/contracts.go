package contracts

import (
	"time"
)

type ServiceRequest struct {
	CipheredRequest       string `json:"ciphered_request,omitempty"`
	CipheredServiceTicket string `json:"ciphered_service_ticket,omitempty"`
}

type ServiceResponse struct {
	CipheredResponse string `json:"ciphered_response,omitempty"`
}

type Response struct {
	Result time.Time `json:"result,omitempty"`
	Nonce  uint64    `json:"nonce,omitempty"`
}

type Request struct {
	ClientID     string        `json:"client_id,omitempty"`
	AccessPeriod time.Duration `json:"access_period,omitempty"`
	Nonce        uint64        `json:"nonce,omitempty"`
}
