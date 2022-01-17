package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/disturb16/finechatbot/broker"
	"github.com/streadway/amqp"
)

// StockCommand is the payload for the stock command.
type StockCommand struct {
	Email      string `json:"email"`
	ChatRoomID int64  `json:"chatroom_id"`
	Message    string `json:"message"`
}

// stockMessage is the incoming payload for the stock command.
type stockMessage struct {
	Type    string        `json:"type"`
	Payload *StockCommand `json:"payload"`
}

type Bot struct {
	messageBroker broker.MessageBroker
}

const (
	apiURL            string = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
	stockCommandTopic string = "stock.command"
)

var ErrInvalidSymbol error = errors.New("invalid stock symbol")

// GetSymbol returns the stock symbol from the message,
// for example if the message is "/stock=googl.us", it returns "googl.us".
func GetSymbol(val string) (string, error) {
	if !strings.HasPrefix(val, "/stock=") || len(val) == 7 {
		return "", ErrInvalidSymbol
	}

	symbol := val[7:]
	return symbol, nil
}

// GetStockDetails returns the share price for the given symbol.
func GetStockDetails(symbol string) (*StockDetails, error) {
	url := fmt.Sprintf(apiURL, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseStockData(data)
}

// New creates a new bot.
func New(b broker.MessageBroker) *Bot {
	return &Bot{
		messageBroker: b,
	}
}

// Listen listens for stock commands and processes them.
func (b *Bot) Listen() error {
	ch, err := b.messageBroker.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = broker.DefaultExchange(ch, stockCommandTopic)
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
		stockCommandTopic, // exchange
		false,
		nil,
	)

	if err != nil {
		return err
	}

	defer ch.QueueUnbind(q.Name, "", stockCommandTopic, nil)

	msgs, err := broker.DefaultConsumer(ch, q)
	if err != nil {
		return err
	}

	foreverChan := make(chan bool)
	go b.listenMessages(msgs)
	<-foreverChan
	return nil
}

func (b *Bot) listenMessages(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Println("Received message:", d.CorrelationId)

		m := &stockMessage{}
		reader := strings.NewReader(string(d.Body))

		err := json.NewDecoder(reader).Decode(m)
		if err != nil {
			log.Println(err)
			continue
		}

		err = b.processStockCommand(m.Payload)
		if err != nil {
			log.Println(err)
		}
	}
}

// processStockCommand processes the stock command.
func (b *Bot) processStockCommand(sc *StockCommand) error {
	exchange := fmt.Sprintf("chatroom.%d", sc.ChatRoomID)

	// Get the stock symbol from the command.
	symbol, err := GetSymbol(sc.Message)
	if err != nil {
		log.Println("Invalid stock symbol:", err)
		return b.sendErrorToUser(exchange, sc.Email, sc.Message)
	}

	// Get the stock price.
	stock, err := GetStockDetails(symbol)
	if err != nil || stock.ClosePrice == "N/D" {
		log.Println("Couldn't get stock price:", err)
		return b.sendErrorToUser(exchange, sc.Email, sc.Message)
	}

	// Send the share price to the chatroom.
	err = b.messageBroker.SendMessage(exchange, exchange, broker.TypeStockRequest, stock.String())
	if err != nil {
		log.Println(err)
		return b.sendErrorToUser(exchange, sc.Email, sc.Message)
	}

	log.Println("Stock price sent to chatroom for symbol:", symbol)
	return nil
}

// sendErrorToUser sends an error message to the user that triggered the event.
func (b *Bot) sendErrorToUser(exchange, email, message string) error {
	key := exchange + "." + email
	payload := "Couldn't process your command: " + message
	return b.messageBroker.SendMessage(exchange, key, broker.TypeCommandError, payload)
}
