package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/helper"
	"github.com/semanggilab/webcore-go/app/logger"
)

type KafkaReceiver interface {
	Run(ctx context.Context, data []byte)
}

// Config menyimpan semua konfigurasi yang dibutuhkan untuk Kafka Consumer.
type KafkaConsumer struct {
	consumer *kafka.Consumer
	config   *config.KafkaConfig
	handler  KafkaReceiver
	topic    string
}

// NewConsumer membuat dan mengembalikan instance kafka.Reader (consumer) baru.
func NewKafkaConsumer(cfg *config.KafkaConfig, h KafkaReceiver) (*KafkaConsumer, error) {
	if cfg == nil || len(cfg.Brokers) == 0 || cfg.Topic == "" || cfg.GroupID == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS, KAFKA_TOPIC, KAFKA_GROUP_ID, dan KAFKA_AUTO_OFFSET_RESET harus di-set")
	}

	kafkaBrokers := strings.Join(cfg.Brokers, ",")

	// Create base config
	config := &kafka.ConfigMap{
		"bootstrap.servers":        kafkaBrokers,
		"group.id":                 cfg.GroupID,
		"auto.offset.reset":        cfg.AutoOffsetReset,
		"socket.timeout.ms":        30000,       // 30 seconds socket timeout
		"socket.keepalive.enable":  true,        // Enable keepalive
		"reconnect.backoff.ms":     10000,       // 10 seconds reconnect backoff
		"reconnect.backoff.max.ms": 300000,      // 5 minutes max backoff
		"statistics.interval.ms":   30000,       // 30 seconds for statistics
		"security.protocol":        "PLAINTEXT", // Use PLAINTEXT as all configurations are now PLAINTEXT
	}

	c, err := kafka.NewConsumer(config)

	if err != nil {
		return nil, err
	}

	// Log broker information for debugging
	logger.Info("Kafka consumer dibuat", "brokers", kafkaBrokers, "group_id", cfg.GroupID)

	// Get and log broker metadata to show advertise listener addresses
	metadata, err := c.GetMetadata(nil, true, 5000)
	if err != nil {
		logger.Error("Gagal mendapatkan metadata Kafka", "error", err)
	} else {
		logger.Info("Metadata Kafka broker", "brokers_count", len(metadata.Brokers))
		for _, broker := range metadata.Brokers {
			logger.Info("Broker advertise listener", "broker_id", broker.ID, "host", broker.Host, "port", broker.Port)
		}
	}

	return &KafkaConsumer{
		consumer: c,
		config:   cfg,
		handler:  h,
		topic:    cfg.Topic,
	}, nil
}

func (kc *KafkaConsumer) Run(ctx context.Context) {
	if kc.handler == nil {
		logger.Warn("Run KafkaConsumer without handler")
		return
	}

	// Track connection state to avoid duplicate error logging
	var lastConnectionError time.Time
	var connectionErrorCount int

	// Subscribe to multiple topics if comma-separated
	topics := strings.Split(kc.topic, ",")
	err := kc.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		logger.Error("Gagal subscribe ke topic", "topic", topics, "error", err)
		return
	}
	logger.Info("Consumer berjalan dan mendengarkan topic", "topics", topics)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Consumer berhenti...")
			return
		default:
			msg, err := kc.consumer.ReadMessage(5 * time.Second) // Timeout agar tidak block selamanya
			if err != nil {
				// Cek jika ini adalah error timeout, yang normal terjadi
				if kerr, ok := err.(kafka.Error); ok && kerr.Code() == kafka.ErrTimedOut {
					continue
				}

				// Handle connection errors with proper logging
				if kerr, ok := err.(kafka.Error); ok {
					switch kerr.Code() {
					case kafka.ErrAllBrokersDown:
						now := time.Now()
						// Only log connection errors every 30 seconds to avoid spam
						if now.Sub(lastConnectionError) > 30*time.Second || connectionErrorCount == 0 {
							logger.Error("Kafka broker tidak tersedia, mencoba reconnect...", "error", err)
							lastConnectionError = now
						}
						connectionErrorCount++
						time.Sleep(5 * time.Second) // Wait before retrying
						continue

					case kafka.ErrUnknownBroker, kafka.ErrInvalidArg, kafka.ErrQueueFull:
						now := time.Now()
						if now.Sub(lastConnectionError) > 60*time.Second || connectionErrorCount == 0 {
							logger.Error("Kafka error terjadi", "error", err, "code", kerr.Code())
							lastConnectionError = now
						}
						connectionErrorCount++
						continue

					default:
						// For other errors, log them but not too frequently
						now := time.Now()
						if now.Sub(lastConnectionError) > 60*time.Second || connectionErrorCount == 0 {
							logger.Error("Kafka error terjadi", "error", err, "code", kerr.Code())
							lastConnectionError = now
						}
						connectionErrorCount++
						continue
					}
				}

				// Handle non-Kafka errors
				now := time.Now()
				if now.Sub(lastConnectionError) > 60*time.Second || connectionErrorCount == 0 {
					logger.Error("Gagal membaca pesan dari Kafka", "error", err)
					lastConnectionError = now
				}
				connectionErrorCount++
				continue
			}

			// Reset connection error count when successfully reading a message
			if connectionErrorCount > 0 {
				logger.Info("Koneksi Kafka berhasil dipulihkan")
				connectionErrorCount = 0
			}

			logger.Debug("Menerima pesan baru", "topic", *msg.TopicPartition.Topic, "offset", msg.TopicPartition.Offset)
			kc.handler.Run(ctx, msg.Value)
		}
	}
}

func (ps *KafkaConsumer) Install(args ...any) error {
	// Tidak melakukan apa-apa
	return nil
}

func (ps *KafkaConsumer) Connect() error {
	// Tidak melakukan apa-apa proses konek hanya dilakukan saat di mode consumer pull message atau publish message di mode producer
	return nil
}

func (ps *KafkaConsumer) Uninstall() error {
	// Tidak melakukan apa-apa
	return nil
}

func (kc *KafkaConsumer) Disconnect() {
	logger.Info("Menutup koneksi Kafka consumer...")
	_ = kc.consumer.Close()
}

// KafkaProducer menyimpan konfigurasi dan instance Kafka producer
type KafkaProducer struct {
	producer *kafka.Producer
	config   *config.KafkaConfig
}

// NewKafkaProducer membuat dan mengembalikan instance KafkaProducer baru
func NewKafkaProducer(cfg *config.KafkaConfig) (*KafkaProducer, error) {
	if cfg == nil || len(cfg.Brokers) == 0 || cfg.Topic == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS dan TOPIC harus di-set")
	}

	kafkaBrokers := strings.Join(cfg.Brokers, ",")

	// Create base config for producer
	config := &kafka.ConfigMap{
		"bootstrap.servers":        kafkaBrokers,
		"socket.timeout.ms":        30000,       // 30 seconds socket timeout
		"socket.keepalive.enable":  true,        // Enable keepalive
		"reconnect.backoff.ms":     10000,       // 10 seconds reconnect backoff
		"reconnect.backoff.max.ms": 300000,      // 5 minutes max backoff
		"statistics.interval.ms":   30000,       // 30 seconds for statistics
		"security.protocol":        "PLAINTEXT", // Use PLAINTEXT as all configurations are now PLAINTEXT
		"acks":                     "all",       // Ensure message durability
		"retries":                  3,           // Number of retries for failed messages
		"linger.ms":                5,           // Batch messages for 5ms
		"batch.size":               16384,       // 16KB batch size
		"compression.type":         "gzip",      // Compress messages to reduce network usage
	}

	p, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat Kafka producer: %v", err)
	}

	// Log broker information for debugging
	logger.Info("Kafka producer dibuat", "brokers", kafkaBrokers, "topic", cfg.Topic)

	// Get and log broker metadata to show advertise listener addresses
	metadata, err := p.GetMetadata(nil, true, 5000)
	if err != nil {
		logger.Error("Gagal mendapatkan metadata Kafka", "error", err)
	} else {
		logger.Info("Metadata Kafka broker", "brokers_count", len(metadata.Brokers))
		for _, broker := range metadata.Brokers {
			logger.Info("Broker advertise listener", "broker_id", broker.ID, "host", broker.Host, "port", broker.Port)
		}
	}

	return &KafkaProducer{
		producer: p,
		config:   cfg,
	}, nil
}

// SendMessage mengirim pesan ke Kafka topic
func (kp *KafkaProducer) SendMessage(ctx context.Context, key string, value []byte) error {
	if kp.producer == nil {
		return fmt.Errorf("producer belum diinisialisasi")
	}

	// Create message
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.config.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
		Headers:        []kafka.Header{}, // Optional headers
	}

	// Channel untuk menerima delivery report
	deliveryChan := make(chan kafka.Event, 10000)

	// Send message
	err := kp.producer.Produce(message, deliveryChan)
	if err != nil {
		return fmt.Errorf("gagal mengirim pesan: %v", err)
	}

	// Wait for delivery report with timeout
	select {
	case event := <-deliveryChan:
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				return fmt.Errorf("pesan gagal terkirim: %v", e.TopicPartition.Error)
			}
			logger.Debug("Pesan berhasil terkirim",
				"topic", *e.TopicPartition.Topic,
				"partition", e.TopicPartition.Partition,
				"offset", e.TopicPartition.Offset,
				"key", key)
		case kafka.Error:
			return fmt.Errorf("Kafka error: %v", e)
		default:
			return fmt.Errorf("event tidak dikenal: %v", e)
		}
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout menunggu delivery report")
	case <-ctx.Done():
		return fmt.Errorf("context dibatalkan")
	}

	return nil
}

// SendJSONMessage mengirim pesan sebagai JSON ke Kafka topic
func (kp *KafkaProducer) SendJSONMessage(ctx context.Context, key string, data interface{}) error {
	jsonData, err := helper.JSONMarshal(data)
	if err != nil {
		return fmt.Errorf("gagal marshal JSON: %v", err)
	}
	return kp.SendMessage(ctx, key, jsonData)
}

// Close menutup koneksi Kafka producer
func (kp *KafkaProducer) Disconnect() {
	if kp.producer != nil {
		logger.Info("Menutup koneksi Kafka producer...")
		// Flush any remaining messages
		flushTimeout := 15 * time.Second
		kp.producer.Flush(int(flushTimeout.Milliseconds()))

		// Close producer
		kp.producer.Close()
		kp.producer = nil
	}
}

// Install implements the interface method
func (kp *KafkaProducer) Install(args ...any) error {
	// Tidak melakukan apa-apa
	return nil
}

// Connect implements the interface method
func (kp *KafkaProducer) Connect() error {
	// Producer tidak memerlukan koneksi eksplisit
	return nil
}

// Uninstall implements the interface method
func (kp *KafkaProducer) Uninstall() error {
	// Tidak melakukan apa-apa
	return nil
}

// Events returns the producer events channel for advanced usage
func (kp *KafkaProducer) Events() <-chan kafka.Event {
	return kp.producer.Events()
}

// Flush flushes any remaining messages in the producer queue
func (kp *KafkaProducer) Flush(timeoutMs int) {
	if kp.producer != nil {
		kp.producer.Flush(timeoutMs)
	}
}
