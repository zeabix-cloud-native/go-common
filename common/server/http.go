package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"net/http"
	"os"
	"os/signal"
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zeabix-cloud-native/go-common/common/config"
	"github.com/zeabix-cloud-native/go-common/common/logs"
	http_reponse "github.com/zeabix-cloud-native/go-common/common/server/http_response"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

func RunHTTPServer(createHandler func(router *gin.Engine) *gin.Engine) {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	RunHTTPServerOnAddr(cfg.AppPort, createHandler)
}

func RunHTTPServerOnAddr(addr string, createHandler func(router *gin.Engine) *gin.Engine) {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	setMiddlewares(router)
	createHandler(router)
	logrus.Info("Starting HTTP server on port : ", addr)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", addr),
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func setMiddlewares(router *gin.Engine) {
	cfg, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}
	router.Use(gin.Recovery())
	router.Use(logs.NewGinStructuredLogger(logrus.StandardLogger()))
	addCorsMiddleware(cfg, router)
}

func addCorsMiddleware(cfg config.Config, router *gin.Engine) {
	allowedOrigins := strings.Split(cfg.CorsOrigins, ";")
	if len(allowedOrigins) == 0 {
		return
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"}, //[]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func Authentication(restrictedApi ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		token, err := parseBearerToken(authorizationHeader)
		if err != nil {
			http_reponse.ErrorResponseHandler(c, 401, err)
			c.Abort()
			return
		}
		userID, access, err := validateToken(token)
		if err != nil {
			http_reponse.ErrorResponseHandler(c, 401, err)
			c.Abort()
			return
		}
		c.Set("USERID", userID)
		if len(restrictedApi) > 0 {
			var count int
			for _, v := range restrictedApi {
				if slices.Contains(access, v) {
					count++
				}
			}
			if count < 1 {
				http_reponse.ErrorResponseHandler(c, 401, http_reponse.ACCESS_DENIED)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func validateToken(tokenstring string) (int, []string, error) {
	cfg, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}
	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return 0, nil, err
	}
	var parsedID interface{}
	var parsedAccess interface{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		parsedID = claims["sub"]
		parsedAccess = claims["aud"]
		if exp, ok := claims["exp"].(float64); ok && float64(time.Now().Unix()) > exp {
			return 0, nil, fmt.Errorf("token expired")
		}
	}
	id, ok := parsedID.(float64)
	if !ok {
		return 0, nil, fmt.Errorf("expected an int value, but got %T", parsedID)
	}
	access, ok := parsedAccess.([]interface{})
	if !ok {
		return 0, nil, fmt.Errorf("expected a slice of strings, but got %T", parsedAccess)
	}
	accessList := make([]string, len(access))
	for i, v := range access {
		if s, ok := v.(string); ok {
			accessList[i] = s
		}
	}
	return int(id), accessList, nil
}

func parseBearerToken(authorizationHeader string) (string, error) {
	const bearerPrefix = "Bearer "
	if len(authorizationHeader) <= len(bearerPrefix) {
		return "", errors.New("invalid authorization header")
	}
	token := authorizationHeader[len(bearerPrefix):]
	return token, nil
}
