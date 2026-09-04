package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nixon-Alexander/resolve_now.git/model"
)

func ReadFile(fileName string) ([]byte, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func UploadFile(
	files []*multipart.FileHeader,
) ([]model.UploadedFile, error) {

	uploadedFiles := make([]model.UploadedFile, 0, len(files))

	uploadDir := filepath.Join(
		"upload",
		"documents",
	)

	// Buat folder kalau tidak ada folder nya
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)

		originalFileName := strings.TrimSuffix(
			filepath.Base(file.Filename),
			fileExt,
		)

		now := time.Now()

		filename := strings.ReplaceAll(
			strings.ToLower(originalFileName),
			" ",
			"-",
		) + "-" +
			fmt.Sprintf("%v", now.Unix()) +
			fileExt

		storagePath := filepath.Join(
			uploadDir,
			filename,
		)

		out, err := os.Create(storagePath)
		if err != nil {
			return nil, err
		}

		readerFile, err := file.Open()
		if err != nil {
			out.Close()
			return nil, err
		}

		_, err = io.Copy(out, readerFile)

		readerFile.Close()
		out.Close()

		if err != nil {
			return nil, err
		}

		uploadedFiles = append(
			uploadedFiles,
			model.UploadedFile{
				FilePath: storagePath,
				FileName: filename,
				MimeType: file.Header.Get("Content-Type"),
			},
		)
	}

	return uploadedFiles, nil
}

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}
