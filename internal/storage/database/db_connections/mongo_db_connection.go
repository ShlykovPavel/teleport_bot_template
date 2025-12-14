package db_connections

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TODO: Добавьте здесь свою БД
// Этот файл содержит пример подключения к MongoDB.
// Вы можете заменить его на PostgreSQL, MySQL или другую БД по необходимости.

// DbConnect устанавливает соединение с MongoDB и возвращает клиент
// TODO: Замените на свою БД если используете PostgreSQL, MySQL и т.д.
func MongoDbConnect(dbUri, username, password string, maxConnectionPoolSize uint64) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(dbUri).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(maxConnectionPoolSize)

	// Добавляем учетные данные, если они указаны
	if username != "" && password != "" {
		clientOpts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Connection error:", err)
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Ping failed:", err)
		return nil, err
	}

	log.Printf("Connected to MongoDB at %s", dbUri)
	fmt.Println("Connection successful! Ready to work.")
	return client, nil
}

// DbDisconnect закрывает соединение с БД
func DbDisconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Println("Disconnect error:", err)
		return err
	}

	fmt.Println("Disconnection successful.")
	return nil
}

// MonitorPool периодически логирует состояние пула соединений с базой данных
func MonitorPool(ctx context.Context, client *mongo.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := client.NumberSessionsInProgress()
			log.Printf("Active DB sessions: %d\n", stats)
		case <-ctx.Done():
			log.Println("Stopping DB pool monitoring")
			return
		}
	}
}
