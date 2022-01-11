package finechatbot

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/disturb16/finechat/broker"
)

const (
	apiURL            string = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
	StockCommandTopic string = "stock_command"
)

var ErrInvalidSymbol error = errors.New("invalid stock symbol")

// StockCommand is the payload for the stock command.
type StockCommand struct {
	Email      string `json:"email"`
	ChatRoomID int64  `json:"chatroom_id"`
	Message    string `json:"message"`
}

// stockMessage is the incoming payload for the stock command.
type stockMessage struct {
	Type    string       `json:"type"`
	Payload StockCommand `json:"payload"`
}

// GetSymbol returns the stock symbol from the message,
// for example if the message is "/stock=googl.us", it returns "googl.us".
func GetSymbol(val string) (string, error) {
	if !strings.HasPrefix(val, "/stock=") || len(val) == 7 {
		return "", ErrInvalidSymbol
	}

	symbol := val[7:]
	return symbol, nil
}

// GetShare returns the share price for the given symbol.
func GetShare(symbol string) (string, error) {
	url := fmt.Sprintf(apiURL, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return parseStockData(resp.Body)
}

// parseStockData parses the stock data from the response body (a csv content).
func parseStockData(data io.Reader) (string, error) {
	csvLines, err := csv.NewReader(data).ReadAll()
	if err != nil {
		return "", err
	}

	closingPrice := csvLines[1][6]

	return closingPrice, nil
}

// Listen listens for stock commands and processes them.
func Listen(b broker.MessageBroker) error {
	ch, err := b.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = broker.DefaultExchange(ch, StockCommandTopic)
	if err != nil {
		return err
	}

	q, err := broker.DefaultQueue(ch, "")
	if err != nil {
		return err
	}

	// Binding for all messages of the chatroom
	err = ch.QueueBind(
		q.Name,            // queue name
		"",                // routing key
		StockCommandTopic, // exchange
		false,
		nil,
	)

	if err != nil {
		return err
	}

	defer ch.QueueUnbind(q.Name, "", StockCommandTopic, nil)

	msgs, err := broker.DefaultConsumer(ch, q)
	if err != nil {
		return err
	}

	foreverChan := make(chan bool)

	go func() {
		for d := range msgs {
			// Send message to websocket corresponding to the chatroom
			m := &stockMessage{}

			err := json.NewDecoder(strings.NewReader(string(d.Body))).Decode(m)
			if err != nil {
				log.Println(err)
				continue
			}

			err = processStockCommand(b, &m.Payload)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	<-foreverChan
	return nil
}

// processStockCommand processes the stock command.
func processStockCommand(b broker.MessageBroker, sc *StockCommand) error {
	exchange := fmt.Sprintf("chatroom.%d", sc.ChatRoomID)

	// Get the stock symbol from the command.
	symbol, err := GetSymbol(sc.Message)
	if err != nil {
		log.Println(err)
		return sendErrorToUser(b, exchange, sc.Email, sc.Message)
	}

	// Get the share price.
	share, err := GetShare(symbol)
	if err != nil || share == "N/D" {
		log.Println(err)
		return sendErrorToUser(b, exchange, sc.Email, sc.Message)
	}

	// Send the share price to the chatroom.
	payload := fmt.Sprintf("%s quote is $%s per share", symbol, share)
	return b.SendMessage(exchange, exchange, broker.TypeStockRequest, payload)
}

// sendErrorToUser sends an error message to the user that triggered the event.
func sendErrorToUser(b broker.MessageBroker, exchange, email, message string) error {
	key := exchange + "." + email
	payload := "Coudln't process your command: " + message
	return b.SendMessage(exchange, key, broker.TypeCommandError, payload)
}
