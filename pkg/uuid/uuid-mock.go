package uuid

import "github.com/stretchr/testify/mock"

type UuidMock struct {
	mock.Mock
}

func (u *UuidMock) NewString() string {
	args := u.Called()
	return args.Get(0).(string)
}
