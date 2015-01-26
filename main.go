package main

import (
	"bytes"
	"fmt"
	twilio "github.com/carlosdp/twiliogo"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	sid := os.Getenv("twilio_sid")
	auth_token := os.Getenv("twilio_auth")
	client := twilio.NewClient(sid, auth_token)

	hand := NewHand()

	fmt.Println(hand)
	msg, err := twilio.NewMessage(client, "2019890712", "2534863751", twilio.Body(hand))

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg.Status)
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
