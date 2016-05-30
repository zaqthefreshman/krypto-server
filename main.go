package main

import (
	"bytes"
	"fmt"
	twilio "github.com/carlosdp/twiliogo"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	client       twilio.Client
	twilioNumber string
)

func main() {
	port := os.Getenv("PORT")
	twilioNumber = os.Getenv("TWILIO_NUMBER")
	sid := os.Getenv("twilio_sid")
	auth_token := os.Getenv("twilio_auth")

	client = twilio.NewClient(sid, auth_token)

	http.HandleFunc("/twilio", handleRequestHand)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func parseTwiloResponse(body io.Reader) url.Values {
	queryStringBytes, _ := ioutil.ReadAll(body)
	queryString := string(queryStringBytes)
	values, err := url.ParseQuery(queryString)
	if err != nil {
		log.Println(err)
	}
	return values
}

func handleRequestHand(w http.ResponseWriter, r *http.Request) {
	hand := NewHand()

	if _, err := w.Write([]byte(hand)); err != nil {
		w.Write(err)
		fmt.Println(err)
	}
}

func NewHand() string {
	deck := make([]int, 0, 52)

	for j := 0; j < 3; j++ {
		for i := 0; i < 11; i++ {
			deck = append(deck, i+1)
		}
	}

	for j := 0; j < 2; j++ {
		for i := 0; i < 6; i++ {
			deck = append(deck, i+12)
		}
	}

	for i := 0; i < 8; i++ {
		deck = append(deck, i+18)
	}

	dest := make([]int, len(deck))
	perm := rand.Perm(len(deck))
	for i, v := range perm {
		dest[v] = deck[i]
	}

	deck = dest

	rand.Seed(int64(time.Now().Nanosecond()))
	offset := 0
	var hand bytes.Buffer
	for i := 0; i < 5; i++ {
		index := rand.Intn(len(deck))
		card := strconv.Itoa(deck[index])
		hand.WriteString(card)
		hand.WriteString(" ")
		deck = append(deck[:index], deck[index+1:]...)

		if i == 1 {
			offset = len(hand.String())
		}
	}

	hand.WriteString("\r\n")
	for i := 0; i < offset; i++ {
		hand.WriteString(" ")
	}

	hand.WriteString(strconv.Itoa(deck[rand.Intn(len(deck))]))

	return hand.String()
}
