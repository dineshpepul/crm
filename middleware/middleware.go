package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	jwksCache jwk.Set
	jwksURL   string
	jwksMu    sync.RWMutex
	once      sync.Once
)

// func JwtAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Get the token from the Authorization header
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
// 			c.Abort()
// 			return
// 		}

// 		err := Token.TokenValid(c)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			c.Abort()
// 			return
// 		}
// 		data, err := Token.ExtractTokenID(c)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			c.Abort()
// 			return
// 		}
// 		// c.AddParam("userId", strconv.FormatUint(uint64(data), 10))
// 		user_id := strconv.FormatUint(uint64(data), 10)
// 		c.Set("userId", user_id)
// 		c.Next()
// 	}
// }

// initialize JWKS once
func initJWKS() error {
	var initErr error
	jwksURL = strings.TrimSuffix(os.Getenv("JWKS_BASE_URL"), "/") + "/.well-known/jwks.json"

	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		set, err := jwk.Fetch(ctx, jwksURL)
		if err != nil {
			initErr = fmt.Errorf("failed to fetch JWKS: %w", err)
			return
		}

		jwksMu.Lock()
		jwksCache = set
		jwksMu.Unlock()
	})

	return initErr
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ensure JWKS is initialized
		jwksMu.RLock()
		cacheEmpty := jwksCache == nil
		jwksMu.RUnlock()

		if cacheEmpty {
			if err := initJWKS(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "JWKS initialization failed",
				})
				return
			}
		}

		// 2. Extract token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Parse and verify
		jwksMu.RLock()
		cache := jwksCache
		jwksMu.RUnlock()

		tok, err := jwt.ParseString(
			tokenStr,
			jwt.WithKeySet(cache),
			jwt.WithValidate(true),
			jwt.WithIssuer("workfast"),
		)

		// 4. Handle JWKS rotation or lookup failure
		if err != nil && strings.Contains(err.Error(), "lookup failed") {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			newSet, fetchErr := jwk.Fetch(ctx, jwksURL)
			if fetchErr == nil {
				jwksMu.Lock()
				jwksCache = newSet
				jwksMu.Unlock()

				tok, err = jwt.ParseString(
					tokenStr,
					jwt.WithKeySet(newSet),
					jwt.WithValidate(true),
					jwt.WithIssuer("workfast"),
				)
			}
		}

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired()) {
				c.AbortWithStatusJSON(440, gin.H{"message": "Token has expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token_invalid", "message": err.Error()})
			return
		}

		// 5. Extract claims safely
		userID, _ := tok.Get("user_id")
		deviceID, _ := tok.Get("device_id")

		c.Set("user_id", userID)
		c.Set("device_id", deviceID)

		c.Next()
	}
}
