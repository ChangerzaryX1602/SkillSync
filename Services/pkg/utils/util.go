package utils

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/valyala/fastjson"
)

var (
	// Session storage
	SessStore *session.Store
	// parser pool
	JsonParserPool *fastjson.ParserPool
)

func init() {
	if JsonParserPool == nil {
		JsonParserPool = new(fastjson.ParserPool)
	}
}
func RandomJitter(ttl time.Duration) time.Duration {
	return ttl + time.Duration(rand.Intn(60))*time.Second
}
