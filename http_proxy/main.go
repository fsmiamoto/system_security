package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	host     = "localhost"
	port     = "3333"
	protocol = "tcp"
)

const errResponse = `HTTP/1.1 403 Forbidden
Content-Type: text/html
Connection: Closed

<html>
    <head>
        <title>Exemplo de resposta HTTP</title>
    </head>
    <body>
        <h1>Acesso não autorizado!</h1>
    </body>
</html>


`

var logger *syslog.Writer

func main() {
	var err error

	logger, err = syslog.Dial("tcp", "192.168.122.103:514",
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "demotag")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen(protocol, host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Listening for connections at %v:%v", host, port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	const bufLen = 1024
	var statusCode int
	defer conn.Close()

	buf := make([]byte, bufLen)

	_, err := conn.Read(buf)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	re := regexp.MustCompile(`(?m)Host:\s+([^\s]+)\s+`)
	matches := re.FindStringSubmatch(string(buf))
	if len(matches) == 0 {
		log.Println("Host header was not found")
		conn.Write([]byte(errResponse))
		return
	}
	host := matches[1]

	clientIp := conn.RemoteAddr().String()
	serverAddr, err := net.ResolveTCPAddr("tcp", host+":80")
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	defer func() {
		msg := fmt.Sprintf("Request - Client: %v - Server: %v - Status: %d\n", clientIp, serverAddr, statusCode)
		log.Printf(msg)
		logger.Info(msg)
	}()

	if strings.Contains(string(buf), "monitorando") {
		log.Println("request denied")
		conn.Write([]byte(errResponse))
		statusCode = 403
		return
	}

	client, err := net.Dial("tcp", matches[1]+":"+"80")
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	_, err = io.Copy(client, bytes.NewBuffer(buf))
	if err != nil {
		log.Println(err)
		return
	}

	respBytes := bytes.NewBuffer(nil)

	io.Copy(respBytes, client)

	statusCodeRe := regexp.MustCompile(`(?m)HTTP/[0-9]\.[0-9] ([0-9]+) .*$`)
	statusMatches := statusCodeRe.FindStringSubmatch(respBytes.String())
	if len(statusMatches) == 0 {
		log.Println("Status code was not found")
		conn.Write([]byte(errResponse))
		return
	}

	statusCode, _ = strconv.Atoi(statusMatches[1])

	io.Copy(conn, respBytes)
}
