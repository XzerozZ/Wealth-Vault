package utils

import (
	"fmt"
	"net/url"
	"strings"

	storage "github.com/supabase-community/storage-go"
)

type StorageClient struct {
	Client *storage.Client
	Bucket string
	URL    string
}

type SupabaseStorage interface {
	Delete(fileURL string) error
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

func (s *StorageClient) Delete(fileURL string) error {
	filePath, err := s.extractPathFromURL(fileURL)
	if err != nil {
		fmt.Printf("⚠️ Invalid Supabase URL: %s\n", fileURL)
		return nil
	}

	_, err = s.Client.RemoveFile(s.Bucket, []string{filePath})
	if err != nil {
		return fmt.Errorf("failed to delete from supabase: %w", err)
	}

	return nil
}

func (s *StorageClient) extractPathFromURL(fullURL string) (string, error) {
	u, err := url.Parse(fullURL)
	if err != nil {
		return "", err
	}

	searchKey := fmt.Sprintf("/%s/", s.Bucket)
	parts := strings.Split(u.Path, searchKey)

	if len(parts) < 2 {
		return "", fmt.Errorf("bucket name not found in URL")
	}

	return parts[1], nil
}
