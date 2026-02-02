package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"

	"golang.org/x/sync/errgroup"
)

type FileStorage interface {
	UploadStream(file io.Reader, filename string, contentType string) (string, error)
}

func UploadBatchFiles(files []*multipart.FileHeader, userID string, folderName string, storage FileStorage) ([]*pb.FileInfo, error) {
	if len(files) == 0 {
		return []*pb.FileInfo{}, nil
	}

	pbFiles := make([]*pb.FileInfo, len(files))
	var g errgroup.Group

	for i, fileHeader := range files {
		index := i
		f := fileHeader

		g.Go(func() error {
			fileData, err := f.Open()
			if err != nil {
				return err
			}
			defer fileData.Close()

			ext := filepath.Ext(f.Filename)
			newFileName := fmt.Sprintf("%s/%s-%d-%d%s", folderName, userID, time.Now().UnixNano(), index, ext)
			contentType := f.Header.Get("Content-Type")

			url, err := storage.UploadStream(fileData, newFileName, contentType)
			if err != nil {
				return err
			}

			pbFiles[index] = &pb.FileInfo{
				Url:      url,
				FileType: contentType,
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return pbFiles, nil
}
