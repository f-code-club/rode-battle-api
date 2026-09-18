package shared

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpService struct {
	url       string
	taskQueue string

	mu      sync.RWMutex
	conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue

	done chan struct{}
}

func NewRabbitService(url string, taskQueue string) (*AmqpService, error) {
	svc := &AmqpService{
		url:       url,
		taskQueue: taskQueue,
		done:      make(chan struct{}),
	}

	if err := svc.connect(); err != nil {
		return nil, err
	}

	go svc.watchAndReconnect()

	return svc, nil
}

func (s *AmqpService) connect() error {
	conn, err := amqp.Dial(s.url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	q, err := ch.QueueDeclare(
		s.taskQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.done:
		_ = ch.Close()
		_ = conn.Close()
		return amqp.ErrClosed
	default:
	}

	s.conn = conn
	s.Channel = ch
	s.Queue = q

	return nil
}

func (s *AmqpService) watchAndReconnect() {
	for {
		select {
		case <-s.done:
			return
		default:
		}

		s.mu.RLock()
		ch := s.Channel
		conn := s.conn
		s.mu.RUnlock()

		if ch == nil || conn == nil {
			if !s.reconnectLoop() {
				return
			}
			continue
		}

		notifyClose := ch.NotifyClose(make(chan *amqp.Error, 1))
		notifyConnClose := conn.NotifyClose(make(chan *amqp.Error, 1))

		select {
		case <-s.done:
			return
		case err, ok := <-notifyClose:
			if !ok {
				log.Println("rabbitmq channel closed — reconnecting...")
			} else {
				log.Printf("rabbitmq channel closed: %v — reconnecting...", err)
			}
		case err, ok := <-notifyConnClose:
			if !ok {
				log.Println("rabbitmq connection closed — reconnecting...")
			} else {
				log.Printf("rabbitmq connection closed: %v — reconnecting...", err)
			}
		}

		if !s.reconnectLoop() {
			return
		}
	}
}

func (s *AmqpService) reconnectLoop() bool {
	s.closeCurrent()

	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-s.done:
			return false
		default:
		}

		log.Println("attempting to reconnect to rabbitmq...")
		if err := s.connect(); err == nil {
			log.Println("successfully reconnected to rabbitmq")
			return true
		} else {
			log.Printf("failed to reconnect to rabbitmq: %v, retrying in %v...", err, backoff)
		}

		select {
		case <-s.done:
			return false
		case <-time.After(backoff):
		}

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}
func (s *AmqpService) closeCurrent() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Channel != nil {
		_ = s.Channel.Close()
		s.Channel = nil
	}
	if s.conn != nil {
		_ = s.conn.Close()
		s.conn = nil
	}
}

func (s *AmqpService) Publish(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	s.mu.RLock()
	ch := s.Channel
	s.mu.RUnlock()

	if ch == nil {
		return amqp.ErrClosed
	}

	return ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

func (s *AmqpService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.done:
		return nil
	default:
		close(s.done)
	}

	var firstErr error
	if s.Channel != nil {
		if err := s.Channel.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		s.Channel = nil
	}
	if s.conn != nil {
		if err := s.conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		s.conn = nil
	}
	return firstErr
}
