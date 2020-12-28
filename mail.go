package mail

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Client struct {
	client mqtt.Client
}

type Mail struct {
	To      []string `json:"to"`
	Title   string   `json:"title"`
	Body    string   `json:"body"`
}

const (
	host              = "tcp://192.168.0.116:1883"
	topic             = "/emails"
	client            = "mail-client"
	autoReconnect     = false
	disconnectTimeout = 5000
)

func NewClient() (*Client, error) {
	// Init Options
	opts := &mqtt.ClientOptions{}
	opts.SetAutoReconnect(autoReconnect)
	opts.SetClientID(clientId())
	opts.AddBroker(host)
	// Init Client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &Client{client}, nil
}

func clientId() string {
	return fmt.Sprintf("%s-%d", client, time.Now().Nanosecond())
}

func (c *Client) NewMail(mail Mail) error {
	// Marshal Struct
	bytes, err := json.Marshal(mail)
	if err != nil {
		return err
	}
	// Publish
	t := c.client.Publish(topic, 0, false, bytes)
	// Wait for Publish
	<- t.Done()
	// Handle Token Error
	if t.Error() != nil {
		return t.Error()
	}
	return nil
}

func (c *Client) Close() {
	c.client.Disconnect(disconnectTimeout)
}