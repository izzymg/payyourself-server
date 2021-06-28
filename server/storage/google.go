package storage

import (
	"context"
	"errors"
	"io"
	"py-server/server"

	"cloud.google.com/go/storage"
)

// GoogleStorer is a UserSaveStorer which uses Google Cloud storage
type GoogleStorer struct {
	client *storage.Client
}

func (gs GoogleStorer) Fetch(userID string) (io.ReadCloser, error) {
	reader, err := gs.client.Bucket("user-saves-1").Object(userID).NewReader(context.TODO())
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, server.ErrNoUserSave
		} else {
			return nil, err
		}
	}
	return reader, nil
}

func (gs GoogleStorer) Close() error {
	return gs.client.Close()
}

func MakeGoogleStorer(ctx context.Context) (*GoogleStorer, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GoogleStorer{
		client,
	}, nil
}
