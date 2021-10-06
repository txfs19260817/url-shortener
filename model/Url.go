package model

import "time"

type Url struct {
	Key         string    `json:"key" bson:"key"`                   // Key is the path part of a short link
	OriginalURL string    `json:"original_url" bson:"original_url"` // OriginalURL is the original URL
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`     // CreatedAt refers to generated time
	ExpireAt    time.Time `json:"expire_at" bson:"expire_at"`       // ExpireAt refers to expiration time
}
