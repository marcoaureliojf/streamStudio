package queue

import (
	"fmt"
	"log"

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func Init(cfg config.Config) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	)

	var err error
	conn, err = amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	log.Println("Connected to RabbitMQ")
}

func GetConnection() *amqp.Connection {
	return conn
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQ(cfg config.Config) (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.DBUser, cfg.DBPassword, cfg.DBHost, 5672))
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir um canal: %w", err)
	}

	q, err := ch.QueueDeclare(
		"schedules", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao declarar a fila: %w", err)
	}
	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}
func (r *RabbitMQ) Publish(body []byte) error {
	err := r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("falha ao publicar mensagem: %w", err)
	}
	log.Printf("Mensagem publicada para a fila %s", r.queue.Name)
	return nil
}
func (r *RabbitMQ) Consume(handler func(d amqp.Delivery)) {
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatalf("Erro ao consumir mensagens: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()
	log.Printf("Aguardando mensagens na fila %s", r.queue.Name)
	<-forever
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
