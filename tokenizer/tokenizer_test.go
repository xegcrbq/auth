package tokenizer

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTokenizer(t *testing.T) {
	tknz := NewTestTokenizer()
	expiration := time.Now().Add(10 * time.Minute)
	cookie, _ := tknz.NewJWTCookie("fingerprint", "sfdhdskgjn", expiration)
	dc, token, err := tknz.ParseDataClaims(cookie.Value)
	assert.Equal(t, dc, &DataClaims{
		Data:           "sfdhdskgjn",
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiration.Unix()},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, token.Raw, cookie.Value)
	assert.Equal(t, token.Claims, dc)
}
