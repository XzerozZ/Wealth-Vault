package helper

import "fmt"

type StorageDeleter interface {
	Delete(url string) error
}

func DeleteFilesAsync(storage StorageDeleter, fileURLs []string) {
	if len(fileURLs) == 0 {
		return
	}

	go func(urls []string) {
		for _, url := range urls {
			if err := storage.Delete(url); err != nil {
				fmt.Printf("⚠️ [AsyncDelete] Failed to delete file %s: %v\n", url, err)
			}
		}
	}(fileURLs)
}
