package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"crypto/rsa"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/crypto/pkcs12"
	"golang.org/x/net/http2"
)

type Notification struct {
	Token   string
	Payload json.RawMessage
}

func main() {

	devP := flag.Bool("d", false, "use development push server")

	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("usage: %s [-d] cert.p12", os.Args[0])
	}

	var server string
	if *devP {
		server = "api.development.push.apple.com"
	} else {
		server = "api.push.apple.com"
	}

	d, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	key, cert, err := pkcs12.Decode(d, "")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key.(*rsa.PrivateKey),
		}},
	}
	tlsConfig.BuildNameToCertificate()

	client := http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tlsConfig,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var notification Notification
		err := json.Unmarshal([]byte(scanner.Text()), &notification)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Post(fmt.Sprintf(
			"https://%s/3/device/%s",
			server,
			url.QueryEscape(notification.Token),
		), "", bytes.NewReader(notification.Payload))

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)

		var rawRes *json.RawMessage
		if len(buf.Bytes()) != 0 {
			rm := json.RawMessage(buf.Bytes())
			rawRes = &rm
		}

		json, err := json.Marshal(struct {
			Status int              `json:"status"`
			Body   *json.RawMessage `json:"body,omitempty"`
		}{resp.StatusCode, rawRes})

		if err != nil {
			log.Fatal(err)
		}

		os.Stdout.Write(json)
		os.Stdout.Write([]byte{'\n'})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
