package kafka

import (
	"fmt"
	"bytes"
	"encoding/json"
	"github.com/Shopify/sarama"
	gobgp "github.com/osrg/gobgp/server"
)

type KafkaClient struct {
	producer sarama.SyncProducer
}

func NewProducer() (*KafkaClient, error) {
	k := &KafkaClient{}
	var err error
	k.producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, nil)

	return k, err
}

func (client *KafkaClient) SendMessage(event_type string, event_body []byte) error {
	value := bytes.NewBuffer(make([]byte, 0, 64))
	value.WriteString(fmt.Sprintf("{\"type\":\"%s\",\"value\":%s}", event_type, event_body))

	msg := &sarama.ProducerMessage{Topic: "edb", Value: sarama.StringEncoder(value.String())}
	_, _, err := client.producer.SendMessage(msg)

	return err
}

func (client *KafkaClient) SendPathListMessage(msg *gobgp.WatchEventBestPath) error {
	for _, path := range msg.PathList {
		bytes, _ := path.MarshalJSON()
		if err :=client.SendMessage("best_path", bytes); err != nil {
			return err
		}
	}

	return nil
}

func (client *KafkaClient) SendPeerStateMessage(msg *gobgp.WatchEventPeerState) error {
	bytes, _ := json.Marshal(msg)
	return client.SendMessage("peer_state", bytes)
}

func (client *KafkaClient) Close() error {
	return client.producer.Close()
}
