package wsclient

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type request struct {
	Id     string            `json:"id"`
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

type bbConfig interface {
	GetBBURL() string
}

type wsClient struct {
	url string
}

func (c *wsClient) Init(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := c.Run(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	err := g.Wait()
	if err != nil {
		log.Printf("Error running WebSockets")
		return err
	}

	return nil
}

func (c *wsClient) Run(ctx context.Context) error {

	header := http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	conn, resp, err := websocket.DefaultDialer.Dial(c.url, header)
	if err != nil {
		log.Printf("RESP %#v", *resp)
		log.Fatalf("dial: %s with status %d", err.Error(), resp.StatusCode)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var id int

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			req := &request{
				Id:     strconv.Itoa(id),
				Method: "ping",
				Params: make(map[string]string),
			}
			id++
			b, err := json.Marshal(req)
			if err != nil {
				return err
			}
			log.Printf("Sending to WS %s", b)
			err = conn.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.Println("write:", err)
				return err
			}
		case <-ctx.Done():
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func NewWSClient(cfg bbConfig) *wsClient {
	wsc := &wsClient{
		url: cfg.GetBBURL(),
	}

	return wsc
}
