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

type StockCommand struct {
	Email      string `json:"email"`
	ChatRoomID int64  `json:"chatroom_id"`
	Message    string `json:"message"`
}

type stockMessage struct {
	Type    string       `json:"type"`
	Payload StockCommand `json:"payload"`
}

func GetSymbol(val string) (string, error) {
	if !strings.HasPrefix(val, "/stock=") || len(val) == 7 {
		return "", ErrInvalidSymbol
	}

	symbol := val[7:]
	return symbol, nil
}

func GetShare(symbol string) (string, error) {
	url := fmt.Sprintf(apiURL, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return parseStockData(resp.Body)
}

func parseStockData(data io.Reader) (string, error) {
	csvLines, err := csv.NewReader(data).ReadAll()
	if err != nil {
		return "", err
	}

	closingPrice := csvLines[1][6]

	return closingPrice, nil
}

func Listen(b *broker.Broker) error {
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

func processStockCommand(b *broker.Broker, sc *StockCommand) error {
	symbol, err := GetSymbol(sc.Message)
	if err != nil {
		return err
	}

	share, err := GetShare(symbol)
	if err != nil {
		return err
	}

	exchange := fmt.Sprintf("chatroom.%d", sc.ChatRoomID)
	payload := fmt.Sprintf("%s quote is $%s per share", symbol, share)

	return b.SendMessage(exchange, exchange, broker.TypeStockRequest, payload)
}
