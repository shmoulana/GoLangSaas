package jwt

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type JWT struct {
	logger     *logrus.Entry
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWT(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, logger *logrus.Entry) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
		logger:     logger,
	}
}

func (j JWT) Create(ttl time.Duration, content interface{}) (result *string, err error) {
	now := time.Now()

	claims := make(jwt.MapClaims)
	claims["dat"] = content
	claims["exp"] = now.Add(ttl).Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.privateKey)
	if err != nil {
		j.logger.Errorf("[CreateToken] failed to create JWT token, err: %+v", err)
		return
	}

	result = &token

	return
}

func (j JWT) Validate(token string) (resp interface{}, err error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			j.logger.Errorf("[ValidateToken] unexpected method: %s, err: %+v", jwtToken.Header["alg"], err)
			return nil, err
		}

		return j.publicKey, nil
	})

	if err != nil {
		return
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		j.logger.Errorf("[ValidateToken] token is invalid")
		return
	}

	resp = claims["dat"]

	return
}
