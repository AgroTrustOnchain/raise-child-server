package walruspkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	walrus_go "github.com/namihq/walrus-go"
)

type walrusProvider struct {
	url       string
	uploadUrl string
	client    *walrus_go.Client
	errLogger *log.Logger
}

type BlobObject struct {
	BlobId string `json:"blobId"`
}

type NewlyCreated struct {
	BlobObject BlobObject `json:"blobObject"`
}

type WalrusResponse struct {
	NewlyCreated *NewlyCreated `json:"newlyCreated,omitempty"`
}

var _walrusProvider *walrusProvider

const (
	walrus_upload_url string = "https://publisher.walrus-testnet.walrus.space/v1/blobs?epochs=5" // default 5 epochs
)

type IWalrusProvider interface {
	UploadImageUrlToWalrus(imageUrl string) string
	FetchBytesImage(blobID string) ([]byte, error)
}

func InitializeWalrusProvider(errLogger *log.Logger) IWalrusProvider {
	if _walrusProvider == nil {
		var walrusClient = walrus_go.NewClient()

		walrusClient.AggregatorURL = []string{
			"https://walrus.space",
			"https://nodes.guru",
		}

		walrusClient.PublisherURL = []string{
			"https://walrus.space",
		}

		_walrusProvider = &walrusProvider{
			url:       "https://aggregator.walrus-testnet.walrus.space/v1/blobs/",
			uploadUrl: walrus_upload_url,
			client:    walrusClient,
			errLogger: errLogger,
		}
	}

	return _walrusProvider
}

// UploadImageUrlToWalrus implements IWalrusProvider.
func (w *walrusProvider) UploadImageUrlToWalrus(imageUrl string) string {
	resp, err := http.Get(imageUrl)
	if err != nil {
		w.errLogger.Println("Error get img url:", err.Error())
		return ""
	}
	defer resp.Body.Close()

	uploadImgReq, err := http.NewRequest(http.MethodPut, w.uploadUrl, resp.Body)
	if err != nil {
		w.errLogger.Println("Error create upload img request:", err.Error())
		return ""
	}
	uploadImgReq.Header.Set("Content-Type", resp.Header.Get("Content-Type"))

	var httpClient = &http.Client{}
	uploadImgRes, err := httpClient.Do(uploadImgReq)
	if err != nil {
		w.errLogger.Println("Error call upload image to Walrus:", err.Error())
		return ""
	}
	defer uploadImgRes.Body.Close()

	body, err := io.ReadAll(uploadImgRes.Body)
	var walrusResult WalrusResponse
	if err := json.Unmarshal(body, &walrusResult); err != nil {
		w.errLogger.Println("Error unmarshal Walrus response:", err.Error())
		return ""
	}

	if walrusResult.NewlyCreated != nil {
		return walrusResult.NewlyCreated.BlobObject.BlobId
	}

	return ""
}

// FetchBytesImage implements IWalrusProvider.
func (w *walrusProvider) FetchBytesImage(blobID string) ([]byte, error) {
	var postUrl string = w.url + blobID

	resp, err := http.Get(postUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch blob: HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
