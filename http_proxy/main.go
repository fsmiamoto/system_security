package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

const (
	host     = "localhost"
	port     = "3333"
	protocol = "tcp"
)

const errResponse = `
<html>
    <head>
        <title>Exemplo de resposta HTTP</title>
    </head>
    <body>
        <h1>Acesso não autorizado!</h1>
    </body>
</html>
`

func main() {
	l, err := net.Listen(protocol, host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Printf("Listening for connections at %v:%v", host, port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	const bufLen = 1024
	defer conn.Close()

	buf := make([]byte, bufLen)

	_, err := conn.Read(buf)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	log.Printf("Request: %s", string(buf))

	if strings.Contains(string(buf), "monitorando") {
		conn.Write([]byte(errResponse))
		return
	}

	re := regexp.MustCompile(`(?m)Host:\s+([^\s]+)\s+`)

	matches := re.FindStringSubmatch(string(buf))
	if len(matches) == 0 {
		log.Println("Host header was not found")
		conn.Write([]byte(errResponse))
		return
	}

	client, err := net.Dial("tcp", matches[1]+":"+"80")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = io.Copy(client, bytes.NewBuffer(buf))
	if err != nil {
		log.Println(err)
		return
	}

	io.Copy(conn, client)
}
