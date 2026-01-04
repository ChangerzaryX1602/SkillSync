package infrastructure

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChangerzaryX1602/SkillSync/internal/datasources"
	"github.com/ChangerzaryX1602/SkillSync/pkg/logs"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"
	"github.com/hibiken/asynq"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/spf13/viper"
)

type Server struct {
	models.Resources
	Version string
	Build   string
	RunEnv  string
	PrdMode bool
}

func NewServer(version, buildTag, runEnv string) (server *Server, err error) {
	// Init server
	server = &Server{
		PrdMode: (viper.GetString("app.env") == "production"),
		Version: version,
		Build:   buildTag,
		RunEnv:  runEnv,
	}

	// connect to DB
	mainDbConn, err := datasources.ConnectDb(datasources.DbConfig{
		DbDriver: "postgres",
		DbName:   viper.GetString("db.postgres.db_name"),
		Host:     viper.GetString("db.postgres.host"),
		Port:     viper.GetInt("db.postgres.port"),
		Username: viper.GetString("db.postgres.username"),
		Password: viper.GetString("db.postgres.password"),
		Timezone: "Asia/Bangkok",
	})
	if err != nil {
		return
	}

	redisStorage, err := datasources.ConnectToRedis(datasources.RedisConfig{
		Host:     viper.GetString("db.redis.host"),
		Port:     viper.GetInt("db.redis.port"),
		Username: viper.GetString("db.redis.username"),
		Password: viper.GetString("db.redis.password"),
		DbName:   viper.GetInt("db.redis.db_name"),
	})
	if err != nil {
		return
	}
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("db.redis.host"),
			viper.GetInt("db.redis.port"),
		),
	}

	asynqClient := asynq.NewClient(redisOpt)

	fastHTTPClient := datasources.NewFastHTTPClient(true)

	// utils.JsonParserPool = new(fastjson.ParserPool)

	jwtResources, err := NewJwt(viper.GetString("jwt.private"))
	if err != nil {
		return
	}

	// init app resources
	server.Resources = NewResources(fastHTTPClient, mainDbConn, nil, redisStorage, asynqClient, jwtResources, viper.GetString("app.embedding_key"))

	// something that use resources place here

	// pre config server
	err = server.configApp()
	if err != nil {
		return
	}

	return
}

func (s *Server) Run() (err error) {
	app := fiber.New(fiber.Config{
		ErrorHandler:      customErrorHandler,
		IdleTimeout:       60 * time.Second,
		ReadBufferSize:    8 * 1024,
		Prefork:           false, // Disabled for Docker
		StreamRequestBody: true,
		// speed up json with goccy/go-json
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		// behide reverse proxy
		EnableTrustedProxyCheck: true,
	})

	// Logger middleware for Fiber that logs HTTP request/response details.
	app.Use(logger.New(s.LogConfig))

	// Recover middleware for Fiber that recovers from panics anywhere in the stack chain and handles the control to the centralized ErrorHandler.
	app.Use(recover.New(recover.Config{EnableStackTrace: !s.PrdMode}))

	// CORS middleware for Fiber that that can be used to enable Cross-Origin Resource Sharing with various options.
	app.Use(cors.New())

	// App Handlers
	s.SetupRoutes(app)

	// Log GO_ENV
	log.Printf("Runtime ENV: %s", viper.GetString("app.env"))
	log.Printf("Version: %s", s.Version)
	log.Printf("Build: %s", s.Build)

	// Listen from a different goroutine

	// Listen HTTP
	go func() {
		if err := app.Listen(":" + viper.GetString("app.port.http")); err != nil {
			log.Panic(err)
		}
	}()

	// Listen HTTPS
	// go func() {
	// 	if err := app.ListenTLS(":"+viper.GetString("app.port.https"), viper.GetString("app.path.cert"), viper.GetString("app.path.priv")); err != nil {
	// 		log.Panic(err)
	// 	}
	// }()

	// Create channel to signify a signal being sent
	quit := make(chan os.Signal, 1)
	// When an interrupt or termination signal is sent, notify the channel
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// This blocks the main thread until an interrupt is received
	<-quit
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	if s.RedisStorage != nil {
		s.RedisStorage.Close()
	}
	if s.AsynqClient != nil {
		s.AsynqClient.Close()
	}
	if s.FastHTTPClient != nil {
		s.FastHTTPClient.CloseIdleConnections()
	}
	fmt.Println("Successful shutdown.")
	return
}

func (s *Server) configApp() (err error) {
	if s.PrdMode {
		s.SessConfig = session.Config{
			Expiration:     8 * time.Hour,
			KeyLookup:      fmt.Sprintf("%s:%s", "cookie", viper.GetString("app.name")),
			CookieSecure:   true,
			CookieHTTPOnly: true,
		}
		logFileWriter := &logs.LogFileWriter{LogPath: "./log/gofiber-skelton", PrintConsole: true}
		s.LogConfig = logger.Config{
			Format:     "[${blue}${time}${reset}] ${status} - ${ip},${ips} ${method} ${host} ${url}\tUserAgent:	${ua}\tReferer: ${referer}\tAuthorization: ${header:Authorization}\tBytesReceived: ${bytesReceived}\tBytesSent: ${bytesSent}\tError: ${red}${error}${reset}\n",
			TimeFormat: "2006-01-02 15:04:05",
			Output:     logFileWriter,
		}
	} else {
		s.SessConfig = session.ConfigDefault
		s.LogConfig = logger.Config{
			Format:     "[${blue}${time}${reset}] ${status} - ${ip},${ips} ${method} ${host} ${url}\nUserAgent:\t${ua}\nReferer:\t${referer}\nAuthorization:\t${header:Authorization}\nBytesReceived:\t${bytesReceived}\nBytesSent:\t${bytesSent}\nError:\t${red}${error}${reset}\n",
			TimeFormat: "2006-01-02 15:04:05",
		}
	}

	// Use redis for session store if available
	if s.RedisStorage != nil {
		s.SessConfig.Storage = s.RedisStorage
	}

	utils.SessStore = session.New(s.SessConfig)

	return
}
