package mongo_db_errors

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// HandleMongoError обрабатывает типичные ошибки MongoDB и возвращает понятное описание
// TODO: Используйте эту функцию в ваших репозиториях для обработки ошибок MongoDB
func HandleMongoError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("document not found")
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return errors.New("database request timeout exceeded")
	}

	if errors.Is(err, context.Canceled) {
		return errors.New("operation canceled")
	}

	var cmdErr mongo.CommandError
	if errors.As(err, &cmdErr) {
		return fmt.Errorf("Command error MongoDB: %s", cmdErr.Message)
	}

	var writeExc mongo.WriteException
	if errors.As(err, &writeExc) {
		if len(writeExc.WriteErrors) > 0 {
			return fmt.Errorf("MongoDB write error: %s", writeExc.WriteErrors[0].Message)
		}
		if writeExc.WriteConcernError != nil {
			return fmt.Errorf("write concern error: %s", writeExc.WriteConcernError.Message)
		}
		return errors.New("unknown MongoDB write error")
	}

	return err
}
