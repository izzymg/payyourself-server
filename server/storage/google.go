package storage

import (
	"context"
	"errors"
	"io"
	"log"
	"py-server/server"

	"cloud.google.com/go/storage"
)

// GoogleStorer is a UserSaveStorer which uses Google Cloud storage
type GoogleStorer struct {
	client *storage.Client
}

func (gs GoogleStorer) Fetch(userID string) (io.ReadCloser, error) {
	log.Printf("fetching google usersave for %s", userID)

	reader, err := gs.client.Bucket("user-saves-1").Object(userID).NewReader(context.TODO())
	if err != nil {
		log.Printf("fetching google usersave failed %s", err)
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, server.ErrNoUserSave
		} else {
			return nil, err
		}
	}
	log.Printf("fetched google usersave for %s", userID)
	return reader, nil
}

func (gs GoogleStorer) Close() error {
	return gs.client.Close()
}

func MakeGoogleStorer(ctx context.Context) (*GoogleStorer, error) {
	log.Printf("creating google storage cloud")
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("created google storage cloud")
	return &GoogleStorer{
		client,
	}, nil
}
