package gcs

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// func (c *GcsClient) UploadFile(key string, file []byte, contentType string, filePath string) error {
func (c *GcsClient) UploadFile(key string, fileHeader *multipart.FileHeader) error {

	ext := filepath.Ext(fileHeader.Filename)

	uniqueFilename := fmt.Sprintf("%s%s", key, ext)
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	signedUrl, err := c.RequestSignedUrl(SignedUrlTypeUpload, c.conf.GcsSignedUrl.BucketName, c.conf.GcsSignedUrl.Path, uniqueFilename)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", signedUrl, file)
	if err != nil {
		return err
	}

	var contentType string
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".svg":
		contentType = "image/svg+xml"
	default:
		return fmt.Errorf("invalid file type")
	}

	req.Header.Add("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
