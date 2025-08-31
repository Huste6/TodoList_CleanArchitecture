package jwt

import (
	"fmt"
	"g09/common"
	"g09/component/tokenprovider"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtProvider struct {
	secret string
	prefix string
}

func NewTokenJWTProvider(secret string, prefix string) *jwtProvider {
	return &jwtProvider{secret: secret, prefix: prefix}
}

type myClaims struct {
	Payload common.TokenPayLoad `json:"payload"`
	jwt.RegisteredClaims
}

type token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token
}

func (j *jwtProvider) SecretKey() string {
	return j.secret
}

func (j *jwtProvider) Generate(data tokenprovider.TokenPayLoad, expiry int) (tokenprovider.Token, error) {
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		Payload: common.TokenPayLoad{
			UID:   data.UserId(),
			URole: data.Role(),
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Local().Add(time.Second * time.Duration(expiry))),
			IssuedAt:  jwt.NewNumericDate(now.Local()),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})

	mytoken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	return &token{
		Token:   mytoken,
		Expiry:  expiry,
		Created: now,
	}, nil
}

func (j *jwtProvider) Validate(tokenString string) (tokenprovider.TokenPayLoad, error) {
	token, err := jwt.ParseWithClaims(tokenString, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}

	if !token.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}

	claims, ok := token.Claims.(*myClaims)
	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}

	return &claims.Payload, nil
}
