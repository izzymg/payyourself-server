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
	bucket *storage.BucketHandle
}

func (gs GoogleStorer) Fetch(userID string) (io.ReadCloser, error) {
	reader, err := gs.bucket.Object(userID).NewReader(context.TODO())
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, server.ErrNoUserSave
		} else {
			return nil, err
		}
	}
	return reader, nil
}

func (gs GoogleStorer) Save(ctx context.Context, userID string) (io.WriteCloser, error) {
	writer := gs.bucket.Object(userID).NewWriter(ctx)
	writer.ObjectAttrs.ContentType = "application/json"
	return writer, nil
}

func (gs GoogleStorer) Close() error {
	return gs.client.Close()
}

func MakeGoogleStorer(ctx context.Context) (*GoogleStorer, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket("user-saves-1")

	return &GoogleStorer{
		client,
		bucket,
	}, nil
}
