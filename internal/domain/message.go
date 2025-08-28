package domain

// Message represents a message consumed from Kafka.
type Message struct {
	Content []byte
	Headers map[string]string
	Key     string
}
