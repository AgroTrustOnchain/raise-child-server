package pkg

import (
	"github.com/stretchr/testify/mock"
)

type walrusMockProvider struct {
	mock.Mock
}

func InitializeWalrusMockProvider() *walrusMockProvider {
	return &walrusMockProvider{}
}

// FetchBytesImage implements walruspkg.IWalrusProvider.
func (w *walrusMockProvider) FetchBytesImage(blobID string) ([]byte, error) {
	var mockData = w.Called(blobID)

	var res1 []byte
	if mockFunc, ok := mockData.Get(0).(func(string) []byte); ok {
		res1 = mockFunc(blobID)
	} else {
		res1 = mockData.Get(0).([]byte)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(0).(func(string) error); ok {
		res2 = mockFunc(blobID)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}
