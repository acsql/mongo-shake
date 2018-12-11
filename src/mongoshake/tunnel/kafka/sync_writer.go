package kafka

import (
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

type SyncWriter struct {
	brokers   []string
	topic     string
	partition int32
	producer  sarama.SyncProducer

	config *Config
}

func NewSyncWriter(address string) (*SyncWriter, error) {
	c := NewConfig()

	topic, brokers, err := parse(address)
	if err != nil {
		return nil, err
	}

	s := &SyncWriter{
		brokers:   brokers,
		topic:     topic,
		partition: defaultPartition,
		config:    c,
	}

	return s, nil
}

func (s *SyncWriter) Start() error {
	producer, err := sarama.NewSyncProducer(s.brokers, s.config.Config)
	if err != nil {
		return err
	}
	s.producer = producer
	return nil
}

func (s *SyncWriter) SimpleWrite(input []byte) error {
	return s.send(input)
}

//func stringHashValue(s string, mod int) int32 {
	//var hashValue int32
	//for _, c := range s {
		//hashValue = 31*hashValue + int32(c)
	//}
	//if hashValue < 0 {
		//return -hashValue
	//}

	//return hashValue % int32(mod)
//}
// producer 按namespace 并行
// // Partition: stringHashValue(namespace, 4),

func (s *SyncWriter) send(input []byte) error {
	// use timestamp as key
	key := strconv.FormatInt(time.Now().UnixNano(), 16)

	msg := &sarama.ProducerMessage{
		Topic:     s.topic,
		Partition: s.partition,
		Key:       sarama.ByteEncoder(key),
		Value:     sarama.ByteEncoder(input),
	}
	_, _, err := s.producer.SendMessage(msg)
	return err
}

func (s *SyncWriter) Close() error {
	return s.producer.Close()
}
