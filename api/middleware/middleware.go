package middleware

import "github.com/gin-gonic/gin"

var (
	Ms []Middleware
)

type Middleware interface {
	Init() ([]gin.HandlerFunc, error)
}

func Middlewares(ms ...Middleware) {
	Ms = append(Ms, ms...)
}
