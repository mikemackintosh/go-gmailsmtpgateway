package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/mail"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"

	"github.com/mhale/smtpd"
)

var gmailService *gmail.Service
var googleConfigFile string

// init
func init() {
	flag.StringVar(&googleConfigFile, "o", "client.json", "Google Oauth Configuration File")
}

// Generate the client
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tok := getTokenFromWeb(config)
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n\nPlease enter code: ", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// mailHandler will handle all received SMTP traffic
func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	var rawMessage = []string{}
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)

	// Add the to and Reply-To
	rawMessage = append(rawMessage, fmt.Sprintf("To: %s\r\n", to[0]))
	rawMessage = append(rawMessage, fmt.Sprintf("Reply-To: %s\r\n", from))
	rawMessage = append(rawMessage, fmt.Sprintf("Return-Path: %s\r\n", from))

	// Loop through the headers and add them in
	for k, header := range msg.Header {
		for _, h := range header {
			rawMessage = append(rawMessage, fmt.Sprintf("%s: %s\r\n", k, h))
		}
	}

	// Add extra linebreak for splitting headers and body
	rawMessage = append(rawMessage, "\r\n")

	// read message body
	buf := new(bytes.Buffer)
	buf.ReadFrom(msg.Body)
	body := buf.String()
	rawMessage = append(rawMessage, body)

	// New message for our gmail service to send
	var message gmail.Message

	// Compose the message
	messageStr := []byte(strings.Join(rawMessage, ""))

	// Place messageStr into message.Raw in base64 encoded format
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	// Send the message
	_, err := gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}
}

func main() {
	flag.Parse()
	// establish context
	ctx := context.Background()

	// read the oauth config
	b, err := ioutil.ReadFile(googleConfigFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/admin-directory_v1-go-quickstart.json
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	// Create a new gmail service using the client
	gmailService, err = gmail.New(client)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	smtpd.ListenAndServe("127.0.0.1:2525", mailHandler, "MyServerApp", "")
}
