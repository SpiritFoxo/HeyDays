package handlers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQManager struct {
	uri           string
	connection    *amqp.Connection
	channel       *amqp.Channel
	isConnected   bool
	mutex         sync.Mutex
	reconnectChan chan struct{}
	closed        bool
}

func NewRabbitMQManager(uri string) *RabbitMQManager {
	manager := &RabbitMQManager{
		uri:           uri,
		isConnected:   false,
		reconnectChan: make(chan struct{}, 1),
		closed:        false,
	}
	go manager.reconnectLoop()
	manager.reconnectChan <- struct{}{}
	return manager
}

func (m *RabbitMQManager) Connect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.isConnected {
		return nil
	}

	var err error
	m.connection, err = amqp.Dial(m.uri)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	m.channel, err = m.connection.Channel()
	if err != nil {
		m.connection.Close()
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	err = m.setupTopology()
	if err != nil {
		m.channel.Close()
		m.connection.Close()
		return fmt.Errorf("failed to set up topology: %v", err)
	}

	go m.monitorConnection(m.connection, m.channel)
	m.isConnected = true
	log.Println("Successfully connected to RabbitMQ")
	return nil
}

func (m *RabbitMQManager) setupTopology() error {
	err := m.channel.ExchangeDeclare(
		"chat_messages",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	_, err = m.channel.QueueDeclare(
		"chat_messages_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = m.channel.QueueBind(
		"chat_messages_queue",
		"chat_messages_queue",
		"chat_messages",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	return nil
}

func (m *RabbitMQManager) monitorConnection(conn *amqp.Connection, ch *amqp.Channel) {
	connCloseChan := conn.NotifyClose(make(chan *amqp.Error, 1))
	chanCloseChan := ch.NotifyClose(make(chan *amqp.Error, 1))

	select {
	case <-connCloseChan:
		log.Println("RabbitMQ connection closed")
	case <-chanCloseChan:
		log.Println("RabbitMQ channel closed")
	}

	m.mutex.Lock()
	m.isConnected = false
	m.mutex.Unlock()

	log.Println("RabbitMQ connection lost, scheduling reconnection...")

	select {
	case m.reconnectChan <- struct{}{}:
	default:
	}
}

func (m *RabbitMQManager) reconnectLoop() {
	for range m.reconnectChan {
		if m.closed {
			return
		}

		if !m.isConnected {
			backoff := 1 * time.Second
			maxBackoff := 30 * time.Second

			for !m.isConnected && !m.closed {
				log.Printf("Attempting to reconnect to RabbitMQ in %v...", backoff)
				time.Sleep(backoff)

				err := m.Connect()
				if err != nil {
					log.Printf("Failed to reconnect: %v", err)
					if backoff.Milliseconds()*2 < maxBackoff.Milliseconds() {
						backoff = backoff * 2
					} else {
						backoff = maxBackoff
					}
				}
			}
		}
	}
}

func (m *RabbitMQManager) GetChannel() (*amqp.Channel, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.isConnected {
		return nil, fmt.Errorf("not connected to RabbitMQ")
	}

	_, err := m.channel.QueueInspect("chat_messages_queue")
	if err != nil {
		m.channel, err = m.connection.Channel()
		if err != nil {
			m.isConnected = false
			select {
			case m.reconnectChan <- struct{}{}:
			default:
			}
			return nil, fmt.Errorf("failed to create a new channel: %v", err)
		}

		err = m.setupTopology()
		if err != nil {
			return nil, fmt.Errorf("failed to set up topology: %v", err)
		}
	}

	return m.channel, nil
}

func (m *RabbitMQManager) PublishMessage(exchange, routingKey string, mandatory, immediate bool, msg amqp.Publishing) error {
	retries := 3
	var err error

	for i := 0; i < retries; i++ {
		var ch *amqp.Channel
		ch, err = m.GetChannel()
		if err != nil {
			log.Printf("Failed to get channel (attempt %d/%d): %v", i+1, retries, err)
			time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
			continue
		}

		err = ch.Publish(exchange, routingKey, mandatory, immediate, msg)
		if err == nil {
			return nil
		}

		log.Printf("Failed to publish message (attempt %d/%d): %v", i+1, retries, err)
		time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
	}

	return fmt.Errorf("failed to publish message after %d attempts: %v", retries, err)
}

func (m *RabbitMQManager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.closed = true

	if m.channel != nil {
		m.channel.Close()
	}

	if m.connection != nil {
		m.connection.Close()
	}

	m.isConnected = false
}
