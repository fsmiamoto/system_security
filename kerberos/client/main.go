package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/crypto"
	service "github.com/fsmiamoto/system_security/kerberos/service/contracts"
	tgs "github.com/fsmiamoto/system_security/kerberos/tgs/contracts"
	"github.com/hokaccha/go-prettyjson"
    "github.com/rs/zerolog/log"
)

var clientID, secretKey, initVector string

var serviceID = "bb1b48f873671ac1bde6fa494b67ba5f"
var accessPeriod = 30 * time.Minute
var httpClient http.Client

var authorizationServer = "http://localhost:3000"
var ticketGrantingServer = "http://localhost:4000"
var serviceServer = "http://localhost:5000"

func main() {
    // Talk to AS
	svcReq := contracts.ServiceRequest{
		AccessPeriod: accessPeriod,
		ServiceID:    serviceID,
		Nonce:        rand.Uint64(),
	}

	data, err := json.Marshal(svcReq)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while marshaling AS request")
    }

	ciphered, err := crypto.Encrypt([]byte(secretKey), []byte(initVector), data)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while encrypting AS request")
    }

	req := contracts.TGTRequest{
		ClientID:               clientID,
		CipheredServiceRequest: hex.EncodeToString(ciphered),
	}

	body, err := json.Marshal(req)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while marshaling TGT request")
    }

	resp, err := httpClient.Post(authorizationServer, "application/json", bytes.NewReader(body))
    if err != nil {
        log.Fatal().Err(err).Msgf("error while marshaling TGT request")
    }

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        log.Fatal().Msgf("error from AS: %v",string(body))
    }

	tgtResp := &contracts.TGTResponse{}

	d := json.NewDecoder(resp.Body)
    if err := d.Decode(tgtResp); err != nil {
        log.Fatal().Err(err).Msg("error while decoding tgt response")
    }

	c, err := hex.DecodeString(tgtResp.CipheredASResponse)
    if err != nil {
        log.Fatal().Err(err).Msg("error while hex decoding AS response")
    }

	u, err := crypto.Decrypt([]byte(secretKey), []byte(initVector), c)
    if err != nil {
        log.Fatal().Err(err).Msg("error while decrypting AS response")
    }

	asResp := &contracts.ASResponse{}

    if err := json.Unmarshal(u, asResp); err != nil {
        log.Fatal().Err(err).Msg("error while unmarshal AS response")
    }

	if asResp.Nonce != svcReq.Nonce {
        log.Fatal().Msgf("error: nonce received does not match %d; want %d", asResp.Nonce, svcReq.Nonce)
	}

	a, _ := prettyjson.Marshal(asResp)
	t, _ := prettyjson.Marshal(tgtResp)

	fmt.Println("Response - AS")
	fmt.Println(string(t))
	fmt.Println(string(a))

    // Talk to TGS

	tgsSvcReq := &tgs.ServiceRequest{
		ClientID:     clientID,
		ServiceID:    serviceID,
		AccessPeriod: accessPeriod,
		Nonce:        rand.Uint64(),
	}

	tgsSvcReqBytes, err := json.Marshal(tgsSvcReq)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while marshaling tgs service request")
    }

	tgsSvcReqBytes, err = crypto.Encrypt([]byte(asResp.KeyClientTGS), []byte(asResp.TGSInitVector), tgsSvcReqBytes)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while encrypting tgs service request")
    }

	str := &tgs.ServiceTicketRequest{
		CipheredServiceRequest: hex.EncodeToString(tgsSvcReqBytes),
		CipheredTGT:            tgtResp.CipheredTGT,
	}

	body, err = json.Marshal(str)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while marshaling service ticket request")
    }

	resp, err = httpClient.Post(ticketGrantingServer, "application/json", bytes.NewReader(body))
    if err != nil {
        log.Fatal().Err(err).Msgf("error while sending request to tgs")
    }

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        log.Fatal().Msgf("error from tgs: %v", string(body))
    }

	tgsRes := &tgs.ServiceTicketResponse{}

	d = json.NewDecoder(resp.Body)
    if err := d.Decode(tgsRes); err != nil {
        log.Fatal().Err(err).Msgf("error while unmarshaling service ticket response")
    }

	tgsResBytes, err := hex.DecodeString(tgsRes.CipheredTGSResponse)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while hex decoding service ticket response")
    }
	tgsResBytes, err = crypto.Decrypt([]byte(asResp.KeyClientTGS), []byte(asResp.TGSInitVector), tgsResBytes)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while decrypting service ticket response")
    }

	tgsResponse := &tgs.TGSResponse{}
    if err := json.Unmarshal(tgsResBytes, tgsResponse); err != nil {
        log.Fatal().Err(err).Msgf("error while unmarshaling tgs response")
    }

	if tgsResponse.Nonce != tgsSvcReq.Nonce {
        log.Fatal().Msgf("error[TGS]: nonce received is different from expected %d; want %d", tgsResponse.Nonce, tgsSvcReq.Nonce)
	}

	fmt.Println("Response - TGS")
	a, _ = prettyjson.Marshal(tgsRes)
	t, _ = prettyjson.Marshal(tgsResponse)

	fmt.Println(string(a))
	fmt.Println(string(t))

    // Talk to service

	r := service.Request{
		ClientID:     clientID,
		AccessPeriod: tgsResponse.AccessPeriod,
	}

	rBytes, err := json.Marshal(r)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while marshaling service request")
    }
	rBytes, err = crypto.Encrypt([]byte(tgsResponse.KeyClientService), []byte(tgsResponse.ServiceInitVector), rBytes)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while encrypting service request")
    }

	svcReqq := &service.ServiceRequest{
		CipheredRequest:       hex.EncodeToString(rBytes),
		CipheredServiceTicket: tgsRes.CipheredServiceTicket,
	}

	svcReqBytes, err := json.Marshal(svcReqq)
    if err != nil {
        log.Fatal().Err(err).Msgf("error while marshaling service request body")
    }

	resp, err = httpClient.Post(serviceServer, "application/json", bytes.NewReader(svcReqBytes))
    if err != nil {
        log.Fatal().Err(err).Msgf("error while sending request to service")
    }
    
    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        log.Fatal().Msgf("error[service]: %v",string(body))
    }

	svcResp := &service.ServiceResponse{}
	d = json.NewDecoder(resp.Body)
    if err := d.Decode(svcResp); err != nil {
        log.Fatal().Err(err).Msgf("error while unmarhsaling service response")
    }

	fmt.Println("Response - Service")
	a, _ = prettyjson.Marshal(svcResp)
	fmt.Println(string(a))

	respBytes, err := hex.DecodeString(svcResp.CipheredResponse)
    if err != nil {
        log.Fatal().Err(err).Msgf("error hex decoding service response")
    }

	respBytes, err = crypto.Decrypt([]byte(tgsResponse.KeyClientService), []byte(tgsResponse.ServiceInitVector), respBytes)
    if err != nil {
        log.Fatal().Err(err).Msgf("error decrypting service response")
    }

	rr := &service.Response{}
    if err := json.Unmarshal(respBytes, rr); err != nil {
        log.Fatal().Err(err).Msgf("error unmarshaling service response")
    }

	a, _ = prettyjson.Marshal(rr)
	fmt.Println(string(a))
}
