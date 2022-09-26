package middleware

import "github.com/shmoulana/Redios/pkg/util/crypt"

type Middleware struct {
	CryptService crypt.CryptService
}
