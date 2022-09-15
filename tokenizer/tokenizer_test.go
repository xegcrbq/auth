package tokenizer

import (
	"fmt"
	"testing"
	"time"
)

func TestTokenizer(t *testing.T) {
	tknz := NewTestTokenizer()
	expiration := time.Now().Add(10 * time.Minute)
	cookie, _ := tknz.NewJWTCookie("fingerprint", "sfdhdskgjn", expiration)
	tknz2 := NewTokenizer([]byte("dhdgvc"))
	fmt.Println(tknz2.ParseDataClaims(cookie.Value))
	fmt.Println(tknz.ParseDataClaims(cookie.Value))
}
