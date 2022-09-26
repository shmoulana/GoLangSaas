package errors

import (
	"errors"

	"github.com/shmoulana/Redios/pkg/dto"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInternalServer = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized")

	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrUnparseableRequestBody    = errors.New("unparseable request body error / invalid request body error")
	ErrAuthTokenExpired          = errors.New("auth token expired")

	ErrUserDataNotFound   = errors.New("user data not found")
	ErrDataNotFound       = errors.New("data not found")
	ErrInsufficientPoints = errors.New("invalid insufficient points")
	ErrInvalidStock       = errors.New("invalid stock is not enough to redeem")
)

var errorMapping = map[error]dto.ErrorResponse{
	ErrBadRequest:                {HTTPCode: 400, Code: 1001, Message: ErrBadRequest.Error()},
	ErrInternalServer:            {HTTPCode: 500, Code: 1002, Message: ErrInternalServer.Error()},
	ErrUnauthorized:              {HTTPCode: 403, Code: 1003, Message: ErrUnauthorized.Error()},
	ErrUnparseableRequestBody:    {HTTPCode: 400, Code: 1010, Message: ErrUnparseableRequestBody.Error()},
	ErrUserDataNotFound:          {HTTPCode: 404, Code: 1004, Message: ErrUserDataNotFound.Error()},
	ErrDataNotFound:              {HTTPCode: 404, Code: 1006, Message: ErrDataNotFound.Error()},
	ErrInvalidUsernameOrPassword: {HTTPCode: 400, Code: 2001, Message: ErrInvalidUsernameOrPassword.Error()},
	ErrAuthTokenExpired:          {HTTPCode: 403, Code: 2002, Message: ErrAuthTokenExpired.Error()},
	ErrInsufficientPoints:        {HTTPCode: 403, Code: 2003, Message: ErrInsufficientPoints.Error()},
	ErrInvalidStock:              {HTTPCode: 403, Code: 2004, Message: ErrInvalidStock.Error()},
}

func GetErrorResponse(er error) (errRes dto.ErrorResponse) {
	errRes, found := errorMapping[er]
	if !found {
		errRes = errorMapping[ErrInternalServer]
	}
	return
}
