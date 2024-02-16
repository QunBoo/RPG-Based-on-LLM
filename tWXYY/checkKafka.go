package tWXYY

import (
	"context"
	"github.com/IBM/sarama"
	"log"
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		sess.MarkMessage(msg, "") // 标记消息为已处理
	}
	return nil
}

func kafkaProduce() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 确保消息被写入所有副本后才认为是成功的
	config.Producer.Retry.Max = 5                    // 最大重试次数
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"myServerIP:9092"}, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: "testTopic",
		Value: sarama.StringEncoder("Hello, World!"),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalln("Failed to send message:", err)
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "your_topic", partition, offset)
}

func kafkaConsumGroup() {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // 确保版本兼容
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 从最早的消息开始消费

	group, err := sarama.NewConsumerGroup([]string{"MyserverIP:9092"}, "your_consumer_group_id", config)
	if err != nil {
		log.Fatalln("Error creating consumer group:", err)
	}
	defer group.Close()

	ctx := context.Background()
	handler := exampleConsumerGroupHandler{}

	// 消费者组循环，确保在消费者出错时可以重新加入
	for {
		if err := group.Consume(ctx, []string{"testTopic"}, handler); err != nil {
			log.Printf("Error from consumer: %v", err)
		}
	}
}
