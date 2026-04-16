package walruspkg

import (
	"log"

	walrus_go "github.com/namihq/walrus-go"
)

type walrusProvider struct {
	client    *walrus_go.Client
	errLogger *log.Logger
}

var _walrusProvider *walrusProvider

type IWalrusProvider interface {
	FetchBytesImage(blobID string) ([]byte, error)
}

func InitializeWalrusProvider(errLogger *log.Logger) IWalrusProvider {
	if _walrusProvider == nil {
		_walrusProvider = &walrusProvider{
			client:    walrus_go.NewClient(),
			errLogger: errLogger,
		}
	}

	return _walrusProvider
}

// FetchBytesImage implements IWalrusProvider.
func (w *walrusProvider) FetchBytesImage(blobID string) ([]byte, error) {
	res, err := w.client.Read(blobID, nil)
	if err != nil {
		w.errLogger.Println(err.Error())
	}

	return res, nil
}
