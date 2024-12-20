package crypto

import (
	"github.com/stretchr/testify/mock"
)

type CryptoMock struct {
	mock.Mock
}

func (c *CryptoMock) Hash(text string) (string, error) {
	args := c.Called(text)
	return args.Get(0).(string), args.Error(1)
}
