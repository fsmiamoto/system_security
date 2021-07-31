package contracts

import "time"

type TGTRequest struct {
	ClientID               string `json:"client_id"`
	CipheredServiceRequest []byte `json:"ciphered_service_request"`
}

type ServiceRequest struct {
	ServiceID    string    `json:"service_id"`
	AccessPeriod time.Time `json:"access_period"`
	Nonce        uint64    `json:"nonce"`
}

type Client struct {
	ID        string
	SecretKey string
}
