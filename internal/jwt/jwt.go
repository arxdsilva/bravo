package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	// this should always be a secret env loaded as []byte
	Secret     string `envconfig:"APP_JWT_SECRET" default:"debug"`
	Expiration int    `envconfig:"APP_JWT_EXPIRATION" default:"720"`
}

type TokenProvider struct {
	expiration int
	secret     string
}

const XAuthenticationToken = "X-Authentication-Token"

func NewProvider(cfg Config) TokenProvider {
	return TokenProvider{
		expiration: cfg.Expiration,
		secret:     cfg.Secret,
	}
}

// New creates a JWT token from the user email
// expires within an hour
func (tp TokenProvider) New(email, userUUID string) (t string, err error) {
	mc := jwt.MapClaims{}
	mc["authorized"] = true
	mc["email"] = email
	mc["user_uuid"] = userUUID
	mc["exp_utc"] = time.Now().Add(time.Duration(tp.expiration) * time.Minute).UTC()
	mc["exp"] = time.Now().Add(time.Duration(tp.expiration) * time.Minute).Unix()
	c := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	return c.SignedString([]byte(tp.secret))
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get(XAuthenticationToken)
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// Validate loads the jwt and ensures that Its valid
func (tp TokenProvider) Validate(r *http.Request) (email string, err error) {
	token, err := tp.verify(r)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("token not valid")
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	email, _ = claims["email"].(string)
	return
}

// ValidateUserUUID loads the jwt and ensures that Its valid
// returns the user uuid
func (tp TokenProvider) ValidateUserUUID(r *http.Request) (uuid string, err error) {
	token, err := tp.verify(r)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("token not valid")
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	uuid, _ = claims["user_uuid"].(string)
	return
}

func (tp TokenProvider) verify(r *http.Request) (t *jwt.Token, err error) {
	tokenString := extractToken(r)
	t, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tp.secret), nil
	})
	if err != nil {
		return t, err
	}
	return t, t.Claims.Valid()
}
