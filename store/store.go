package store

import (
	"context"
	"fmt"
	"math/rand"
	"server/objects"
	"time"
)

// IPostStore is the database interface for storing Posts
type IPostStore interface {
	Get(ctx context.Context, in *objects.GetRequest) (*objects.Post, error)
	List(ctx context.Context, in *objects.ListRequest) ([]*objects.Post, error)
	Create(ctx context.Context, in *objects.CreateRequest) error
	UpdateViewed(ctx context.Context, in *objects.UpdateRequest) error
	Delete(ctx context.Context, in *objects.DeleteRequest) error
}

func init() {
	rand.Seed(time.Now().UTC().Unix())
}

// GenerateUniqueID will returns a time based sortable unique id
func GenerateUniqueID() string {
	word := []byte("0987654321")
	rand.Shuffle(len(word), func(i, j int) {
		word[i], word[j] = word[j], word[i]
	})
	now := time.Now().UTC()
	return fmt.Sprintf("%010v-%010v-%s", now.Unix(), now.Nanosecond(), string(word))
}
