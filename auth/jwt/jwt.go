package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

func Generate(secretKey string, payloads jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, payloads).SignedString([]byte(secretKey))
}

func Parse(req *http.Request, secret string, out jwt.Claims) (string, error) {
	var parser = &jwt.Parser{UseJSONNumber: true}
	token, err := ParseFromRequest(req, secret, request.WithParser(parser), request.WithClaims(out))
	if err != nil {
		return "", err
	}
	return token.Raw, nil
}

func ParseFromRequest(req *http.Request, secret string, opts ...request.ParseFromRequestOption) (*jwt.Token, error) {
	return request.ParseFromRequest(
		req,
		request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil },
		opts...
	)
}
