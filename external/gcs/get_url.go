package gcs

func (c *GcsClient) getUrl(signedUrlType SignedUrlType) string {
	url := c.conf.GcsSignedUrl.DownloadUrl
	if signedUrlType == SignedUrlTypeUpload {
		url = c.conf.GcsSignedUrl.UploadUrl
	}
	return url
}
