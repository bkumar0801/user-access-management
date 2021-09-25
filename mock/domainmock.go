package mock

import (
	"errors"
)

//ErrMock ...
type ErrMock int

const (
	//DBOperationError ...
	DBOperationError ErrMock = iota

	//DBConnectionError ...
	DBConnectionError

	//DBConnectionTimeout ...
	DBConnectionTimeout

	//DBDuplicateEntry ...
	DBDuplicateEntry

	//DBNoEntry ...
	DBNoEntry

	//DBInvalidOperation ...
	DBInvalidOperation

	//DBInvalidDateTime ...
	DBInvalidDateTime

	//OK ...
	OK
)

//ServiceMock ...
type ServiceMock struct {
	Err ErrMock
}

//DBStatus is a mock implementation of DBStatus
func (s *ServiceMock) DBStatus() (bool, error) {
	if s.Err == DBConnectionError {
		return false, errors.New("mock DB connection error")
	}

	if s.Err == DBConnectionTimeout {
		return false, nil
	}

	return true, nil
}
