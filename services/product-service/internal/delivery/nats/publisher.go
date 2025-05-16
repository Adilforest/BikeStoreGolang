package nats

import (
 "encoding/json"

 "github.com/nats-io/nats.go"
)

type OrderProcessedEvent struct {
 OrderID string `json:"order_id"`
 Status  string `json:"status"`
 Message string `json:"message"`
}

type Publisher interface {
 PublishOrderProcessed(event OrderProcessedEvent) error
}

type natsPublisher struct {
 nc *nats.Conn
}

func NewPublisher(nc *nats.Conn) Publisher {
 return &natsPublisher{nc: nc}
}

func (p *natsPublisher) PublishOrderProcessed(event OrderProcessedEvent) error {
 data, err := json.Marshal(event)
 if err != nil {
  return err
 }
 return p.nc.Publish("order.processed", data)
}
