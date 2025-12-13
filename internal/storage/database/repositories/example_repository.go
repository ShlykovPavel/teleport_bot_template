package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Добавьте здесь свою MongoDB коллекцию и репозиторий
// Это пример структуры для демонстрации паттерна Repository

// ExampleDocument пример структуры документа для MongoDB
// TODO: Замените на вашу модель данных
type ExampleDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Data      string             `bson:"data"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// ExampleRepository интерфейс для работы с примером коллекции
// TODO: Создайте свой интерфейс репозитория
type ExampleRepository interface {
	Create(ctx context.Context, doc *ExampleDocument) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*ExampleDocument, error)
	Update(ctx context.Context, doc *ExampleDocument) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// TODO: Реализуйте интерфейс репозитория
// Пример:
// type exampleRepository struct {
//     collection *mongo.Collection
//     logger     *slog.Logger
// }
//
// func NewExampleRepository(db *mongo.Database, logger *slog.Logger) ExampleRepository {
//     return &exampleRepository{
//         collection: db.Collection("examples"),
//         logger:     logger,
//     }
// }
//
// func (r *exampleRepository) Create(ctx context.Context, doc *ExampleDocument) error {
//     doc.CreatedAt = time.Now()
//     doc.UpdatedAt = time.Now()
//     result, err := r.collection.InsertOne(ctx, doc)
//     if err != nil {
//         return err
//     }
//     doc.ID = result.InsertedID.(primitive.ObjectID)
//     return nil
// }
//
// ... остальные методы

// TODO: Используйте db_errors для обработки ошибок MongoDB
// Пример из основного проекта:
// import "your-module/internal/storage/database/db_errors"
// if err != nil {
//     return db_errors.HandleMongoError(err, "operation description")
// }
