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
	"github.com/hokaccha/go-prettyjson"
)

var clientID, secretKey, initVector string

var httpClient http.Client

var authorizationServer = "http://localhost:3000"

func main() {
	svcReq := contracts.ServiceRequest{
		AccessPeriod: 30 * time.Minute,
		ServiceID:    "1234a",
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

	a, _ := prettyjson.Marshal(asResp)
	t, _ := prettyjson.Marshal(tgtResp)

	fmt.Println("Response - AS")
	fmt.Println("TGT:", string(t))
	fmt.Println("AS:", string(a))
}

// helper for dumb err checking, don't do this on production
func mustBeNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
