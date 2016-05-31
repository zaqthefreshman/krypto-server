package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	twilio "github.com/carlosdp/twiliogo"
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

	http.HandleFunc("/", handleWeb)
	http.HandleFunc("/twilio", handleTwilio)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleWeb(w http.ResponseWriter, r *http.Request) {
	hand := NewHand()
	fmtedHand, _ := fmtRow(hand.Row)

	t, err := template.ParseFiles("./static/hand.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	respHand := struct {
		Row    string
		Target string
	}{
		fmtedHand,
		strconv.Itoa(hand.Target),
	}
	t.Execute(w, respHand)
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

func handleTwilio(w http.ResponseWriter, r *http.Request) {
	hand := NewHand()
	fmtedHand, offset := fmtRow(hand.Row)

	if _, err := w.Write([]byte(fmtedHand + "\r\n" + fmtTarget(offset) + strconv.Itoa(hand.Target))); err != nil {
		fmt.Println(err)
	}
}

type Hand struct {
	Row    []int
	Target int
}

func fmtRow(row []int) (string, int) {
	offset := 0
	var hand bytes.Buffer
	for i, card := range row {
		hand.WriteString(strconv.Itoa(card))

		if i != 4 {
			hand.WriteString(", ")
		}

		if i == 1 {
			offset = len(hand.String())
		}
	}

	return hand.String(), offset
}

func fmtTarget(offset int) string {
	var spaces bytes.Buffer
	for i := 0; i < offset; i++ {
		spaces.WriteString(" ")
	}

	return spaces.String()
}

func NewHand() Hand {
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

	var hand Hand
	for i := 0; i < 5; i++ {
		index := rand.Intn(len(deck))
		card := deck[index]

		// Remove card from deck
		deck = append(deck[:index], deck[index+1:]...)

		hand.Row = append(hand.Row, card)
	}

	index := rand.Intn(len(deck))
	card := deck[index]

	// Remove card from deck
	deck = append(deck[:index], deck[index+1:]...)

	hand.Target = card

	return hand
}
