package utils

import (
	"github.com/IBM/sarama"
	"log"
)

type ExampleConsumerGroupHandler struct{}

func (ExampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ExampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h ExampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		sess.MarkMessage(msg, "") // 标记消息为已处理
	}
	return nil
}
