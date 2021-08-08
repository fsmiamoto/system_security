package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/fsmiamoto/system_security/kerberos/as/contracts"
	"github.com/fsmiamoto/system_security/kerberos/crypto"
	tgs "github.com/fsmiamoto/system_security/kerberos/tgs/contracts"
	"github.com/hokaccha/go-prettyjson"
)

var clientID, secretKey, initVector string

var serviceID = "bb1b48f873671ac1bde6fa494b67ba5f"
var accessPeriod = 30 * time.Minute
var httpClient http.Client

var authorizationServer = "http://localhost:3000"
var ticketGrantingServer = "http://localhost:4000"

func main() {
	svcReq := contracts.ServiceRequest{
		AccessPeriod: accessPeriod,
		ServiceID:    serviceID,
		Nonce:        rand.Uint64(),
	}

	data, err := json.Marshal(svcReq)
	mustBeNil(err)

	ciphered, err := crypto.Encrypt([]byte(secretKey), []byte(initVector), data)
	mustBeNil(err)

	req := contracts.TGTRequest{
		ClientID:               clientID,
		CipheredServiceRequest: hex.EncodeToString(ciphered),
	}

	body, err := json.Marshal(req)
	mustBeNil(err)

	resp, err := httpClient.Post(authorizationServer, "application/json", bytes.NewReader(body))
	mustBeNil(err)

	tgtResp := &contracts.TGTResponse{}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(tgtResp)
	mustBeNil(err)

	c, err := hex.DecodeString(tgtResp.CipheredASResponse)
	mustBeNil(err)

	u, err := crypto.Decrypt([]byte(secretKey), []byte(initVector), c)
	mustBeNil(err)

	asResp := &contracts.ASResponse{}

	mustBeNil(json.Unmarshal(u, asResp))

	if asResp.Nonce != svcReq.Nonce {
		log.Fatalf("Nonce received does not match %d; want %d", asResp.Nonce, svcReq.Nonce)
	}

	a, _ := prettyjson.Marshal(asResp)
	t, _ := prettyjson.Marshal(tgtResp)

	fmt.Println("Response - AS")
	fmt.Println(string(t))
	fmt.Println(string(a))

	tgsSvcReq := &tgs.ServiceRequest{
		ClientID:     clientID,
		ServiceID:    serviceID,
		AccessPeriod: accessPeriod,
		Nonce:        rand.Uint64(),
	}

	tgsSvcReqBytes, err := json.Marshal(tgsSvcReq)
	mustBeNil(err)

	tgsSvcReqBytes, err = crypto.Encrypt([]byte(asResp.KeyClientTGS), []byte(asResp.TGSInitVector), tgsSvcReqBytes)
	mustBeNil(err)

	str := &tgs.ServiceTicketRequest{
		CipheredServiceRequest: hex.EncodeToString(tgsSvcReqBytes),
		CipheredTGT:            tgtResp.CipheredTGT,
	}

	body, err = json.Marshal(str)
	mustBeNil(err)

	resp, err = httpClient.Post(ticketGrantingServer, "application/json", bytes.NewReader(body))
	mustBeNil(err)

	tgsRes := &tgs.ServiceTicketResponse{}

	d = json.NewDecoder(resp.Body)
	mustBeNil(d.Decode(tgsRes))

	tgsResBytes, err := hex.DecodeString(tgsRes.CipheredTGSResponse)
	tgsResBytes, err = crypto.Decrypt([]byte(asResp.KeyClientTGS), []byte(asResp.TGSInitVector), tgsResBytes)

	tgsResponse := &tgs.TGSResponse{}
	mustBeNil(json.Unmarshal(tgsResBytes, tgsResponse))

	if tgsResponse.Nonce != tgsSvcReq.Nonce {
		log.Fatalf("Nonce is different from expected %d; want %d", tgsResponse.Nonce, tgsSvcReq.Nonce)
	}

	fmt.Println("Response - TGS")
	a, _ = prettyjson.Marshal(tgsRes)
	t, _ = prettyjson.Marshal(tgsResponse)

	fmt.Println(string(a))
	fmt.Println(string(t))
}

// helper for dumb err checking, don't do this on production
func mustBeNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
