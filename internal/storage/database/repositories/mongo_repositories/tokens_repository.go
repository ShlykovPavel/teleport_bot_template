package mongo_repositories

import (
	"context"
	mongo_db_errors "template-external-api-service/internal/storage/database/db_errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AccessToken - access token structure in database
type AccessToken struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID int                `bson:"user_id"`
	Token  string             `bson:"token"`
}

type TokensRepository interface {
	SaveAccessToken(ctx context.Context, userID int, token string) error
	DeleteAccessToken(ctx context.Context, token string) error
	CheckAccessToken(ctx context.Context, token string) (bool, error)
	GetAccessTokenInfo(ctx context.Context, token string) (*AccessToken, error)
	UpdateAccessToken(ctx context.Context, userID int, token string) error
}

// TokenRepository - repository for working with access tokens
type TokenRepositoryCollection struct {
	tokensCol *mongo.Collection
}

func NewTokenRepositoryImpl(db *mongo.Database) *TokenRepositoryCollection {
	return &TokenRepositoryCollection{
		tokensCol: db.Collection("access_tokens"),
	}
}

// SaveAccessToken saves access token on login
func (r *TokenRepositoryCollection) SaveAccessToken(ctx context.Context, userID int, token string) error {
	// First delete all old tokens for this user (to have only one active token)
	_, err := r.tokensCol.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return mongo_db_errors.HandleMongoError(err)
	}

	// Save new token
	t := AccessToken{
		UserID: userID,
		Token:  token,
	}
	_, err = r.tokensCol.InsertOne(ctx, t)
	return mongo_db_errors.HandleMongoError(err)
}

func (r *TokenRepositoryCollection) UpdateAccessToken(ctx context.Context, userID int, token string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"token": token}}

	_, err := r.tokensCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return mongo_db_errors.HandleMongoError(err)
	}
	return nil
}

// DeleteAccessToken deletes access token on logout
func (r *TokenRepositoryCollection) DeleteAccessToken(ctx context.Context, token string) error {
	_, err := r.tokensCol.DeleteOne(ctx, bson.M{"token": token})
	return mongo_db_errors.HandleMongoError(err)
}

// CheckAccessToken validates access token
func (r *TokenRepositoryCollection) CheckAccessToken(ctx context.Context, token string) (bool, error) {
	var t AccessToken
	err := r.tokensCol.FindOne(ctx, bson.M{"token": token}).Decode(&t)
	if err == mongo.ErrNoDocuments {
		return false, mongo.ErrNoDocuments
	}
	if err != nil {
		return false, mongo_db_errors.HandleMongoError(err)
	}
	return true, nil
}

// GetAccessTokenInfo returns token information (for access operations)
func (r *TokenRepositoryCollection) GetAccessTokenInfo(ctx context.Context, token string) (*AccessToken, error) {
	var t AccessToken
	err := r.tokensCol.FindOne(ctx, bson.M{"token": token}).Decode(&t)
	if err != nil {
		return nil, mongo_db_errors.HandleMongoError(err)
	}
	return &t, nil
}
