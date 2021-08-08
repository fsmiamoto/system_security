package contracts

import "time"

type ServiceTicketRequest struct {
	CipheredServiceRequest string `json:"ciphered_service_request"`
	CipheredTGT            string `json:"ciphered_tgt"`
}

type ServiceTicketResponse struct {
	CipheredTGSResponse   string `json:"ciphered_tgs_response"`
	CipheredServiceTicket string `json:"ciphered_service_ticket"`
}

type TGSResponse struct {
	KeyClientService string        `json:"key_client_service,omitempty"`
	AccessPeriod     time.Duration `json:"access_period,omitempty"`
	Nonce            uint64        `json:"nonce,omitempty"`
}

type ServiceTicket struct {
	ClientID         string        `json:"client_id"`
	AccessPeriod     time.Duration `json:"access_period"`
	KeyClientService string        `json:"key_client_service"`
}

type ServiceRequest struct {
	ClientID     string        `json:"client_id"`
	ServiceID    string        `json:"service_id"`
	AccessPeriod time.Duration `json:"access_period"`
	Nonce        uint64        `json:"nonce"`
}

type Service struct {
	ID         string `json:"id"`
	SecretKey  []byte `json:"secret_key"`
	InitVector []byte `json:"init_vector"`
}
