package nats

import (
 "encoding/json"
 "log"

 "github.com/nats-io/nats.go"
)

type OrderCreatedEvent struct {
 OrderID string `json:"order_id"`
 UserID  string `json:"user_id"`
 Items   []struct {
  ProductID string  `json:"product_id"`
  Quantity  int32   `json:"quantity"`
  Price     float64 `json:"price"`
 } `json:"items"`
 Total   float64 `json:"total"`
 Address string  `json:"address"`
 Status  string  `json:"status"`
}

func SubscribeOrderCreated(nc *nats.Conn, handle func(OrderCreatedEvent)) error {
 _, err := nc.Subscribe("order.created", func(m *nats.Msg) {
  var event OrderCreatedEvent
  if err := json.Unmarshal(m.Data, &event); err != nil {
   log.Printf("Failed to unmarshal order.created event: %v", err)
   return
  }
  handle(event)
 })
 return err
}
