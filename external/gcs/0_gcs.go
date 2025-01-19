package gcs

import (
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

type (
	SignedUrlType string
	SignedUrlRes  struct {
		Status     string              `json:"status"`
		Resource   string              `json:"resource"`
		Signature  string              `json:"signature"`
		Expiration SignedUrlExpiration `json:"expiration"`
		SignedUrl  string              `json:"signedurl"`
		Error      struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}
	SignedUrlExpiration struct {
		Seconds  int       `json:"seconds"`
		Relative int       `json:"relative"`
		IsoTime  time.Time `json:"ISO"`
	}
	cachedUrl struct {
		url       string
		expiresAt time.Time
	}
)

const (
	SignedUrlTypeUpload   SignedUrlType = "U"
	SignedUrlTypeDownload SignedUrlType = "D"
)

type (
	GcsClientService interface {
		RequestSignedUrl(signedType SignedUrlType, bucketName string, bucketPath string, filename string) (string, error)
		Download(filename string) (string, error)
		UploadFile(key string, fileHeader *multipart.FileHeader) error
	}

	GcsClient struct {
		conf       *config.Config
		cache      map[string]cachedUrl
		cacheMutex sync.Mutex
		client     *http.Client
	}
)

func NewGcsClient(conf *config.Config, client *http.Client) GcsClientService {
	if client == nil {
		client = &http.Client{}
	}

	return &GcsClient{
		conf:   conf,
		cache:  make(map[string]cachedUrl),
		client: client, // Add the client to the adapter
	}
}
