package utils

import (
	"fmt"
	"io"

	storage "github.com/supabase-community/storage-go"
)

type StorageClient struct {
	Client *storage.Client
	Bucket string
	URL    string
}

func NewStorageClient(url, key, bucket string) (*StorageClient, error) {
	if url == "" || key == "" || bucket == "" {
		return nil, fmt.Errorf("missing supabase config")
	}

	client := storage.NewClient(url, key, nil)
	if client == nil {
		return nil, fmt.Errorf("failed to create storage client: invalid Supabase configuration")
	}

	return &StorageClient{
		Client: client,
		Bucket: bucket,
		URL:    url,
	}, nil
}

func (s *StorageClient) UploadStream(file io.Reader, filename string, contentType string) (string, error) {
	opts := storage.FileOptions{
		ContentType: &contentType,
		Upsert:      func() *bool { b := true; return &b }(),
	}

	_, err := s.Client.UploadFile(s.Bucket, filename, file, opts)
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	publicURL := fmt.Sprintf("%s/object/public/%s/%s", s.URL, s.Bucket, filename)

	return publicURL, nil
}
