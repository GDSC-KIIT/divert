package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// URLShorten defines structure for stored urls
type URLShorten struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OriginalURL      string             `json:"original_url,omitempty" bson:"original_url"`
	ShortenedURLCode string             `json:"shortened_url_code,omitempty" bson:"shortened_url_code"`
	ClickCount       int                `json:"click_count" bson:"click_count"`
}
