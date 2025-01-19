package gcs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *GcsClient) RequestSignedUrl(signedType SignedUrlType, bucketName string, bucketPath string, filename string) (string, error) {
	signedUrl := c.getUrl(signedType)
	formData := url.Values{
		"bucket":   {bucketName},
		"path":     {bucketPath},
		"filename": {filename},
		"expiry":   {strconv.Itoa(c.conf.GcsSignedUrl.Expired)},
	}
	req, err := http.NewRequest(http.MethodPost, signedUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("x-apikey", c.conf.GcsSignedUrl.ApiKey)
	req.Header.Set("x-service-name", c.conf.GcsSignedUrl.ServiceName)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	signedUrlsResp := SignedUrlRes{}
	err = json.Unmarshal(responseBody, &signedUrlsResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %s", string(responseBody))
	}

	if signedUrlsResp.Error.Code != 0 {
		return "", fmt.Errorf("error code %d: %s", signedUrlsResp.Error.Code, signedUrlsResp.Error.Message)
	}

	return signedUrlsResp.SignedUrl, nil
}
