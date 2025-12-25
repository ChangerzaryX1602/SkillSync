package infrastructure

import (
	"fmt"

	"github.com/ChangerzaryX1602/SkillSync/internal/datasources"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"github.com/gofiber/storage/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func NewResources(fasthttpClient *fasthttp.Client, mainDbConn *gorm.DB, logDbConn *mongo.Database, redisStorage *redis.Storage, jwtResources *models.JwtResources, embbedingKey string) models.Resources {
	return models.Resources{
		FastHTTPClient: fasthttpClient,
		MainDbConn:     mainDbConn,
		LogDbConn:      logDbConn,
		RedisStorage:   redisStorage,
		JwtResources:   jwtResources,
		EmbeddingKey:   embbedingKey,
	}
}

func NewJwt(privateKeyPath string) (jwtResources *models.JwtResources, err error) {
	jwtResources = new(models.JwtResources)
	jwtResources.JwtSignKey, jwtResources.JwtVerifyKey, jwtResources.JwtSigningMethod, err = datasources.NewJwtLocalKey(privateKeyPath)
	jwtResources.JwtKeyfunc = func(token *jwt.Token) (publicKey interface{}, err error) {
		if jwtResources.JwtVerifyKey == nil {
			err = fmt.Errorf("JWTVerifyKey not init yet")
		}
		return jwtResources.JwtVerifyKey, err
	}
	jwtResources.JwtParser = jwt.NewParser()
	return
}
