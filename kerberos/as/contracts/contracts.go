package contracts

import "time"

type TGTRequest struct {
	ClientID               string `json:"client_id"`
	CipheredServiceRequest string `json:"ciphered_service_request"`
}

type TGTResponse struct {
	CipheredASResponse string `json:"ciphered_as_response"`
	CipheredTGT        string `json:"ciphered_tgt"`
}

type ASResponse struct {
	KeyClientTGS  string `json:"key_client_tgs"`
	TGSInitVector string `json:"tgs_init_vector"`
	Nonce         uint64 `json:"nonce"`
}

type TGT struct {
	ClientID     string        `json:"client_id"`
	CreatedAt    time.Time     `json:"created_at"`
	AccessPeriod time.Duration `json:"access_period"`
	KeyClientTGS string        `json:"key_client_tgs"`
}

type ServiceRequest struct {
	ServiceID    string        `json:"service_id"`
	AccessPeriod time.Duration `json:"access_period"`
	Nonce        uint64        `json:"nonce"`
}

type Client struct {
	ID         string `json:"id"`
	SecretKey  []byte `json:"secret_key"`
	InitVector []byte `json:"init_vector"`
}
