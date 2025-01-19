package gcs

import (
	"fmt"
	"time"
)

func (c *GcsClient) Download(filename string) (string, error) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	cacheKey := fmt.Sprintf("%s/%s", c.conf.GcsSignedUrl.Path, filename)
	if cached, found := c.cache[cacheKey]; found && cached.expiresAt.After(time.Now()) {
		return cached.url, nil
	}

	signedUrl, err := c.RequestSignedUrl(SignedUrlTypeDownload, c.conf.GcsSignedUrl.BucketName, c.conf.GcsSignedUrl.Path, filename)
	if err != nil {
		return signedUrl, err
	}

	c.cache[cacheKey] = cachedUrl{
		url:       signedUrl,
		expiresAt: time.Now().Add(1 * time.Hour),
	}
	return signedUrl, nil
}
