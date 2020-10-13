package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// URLShorten defines structure for stored urls
type URLShorten struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OriginalURL      string             `json:"originalUrl,omitempty"`
	ShortenedURLCode string             `json:"shortenedUrlCode,omitempty"`
	ClickCount       int                `json:"clickCount"`
}
