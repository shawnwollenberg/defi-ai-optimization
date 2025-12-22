module github.com/defioptimization/api

go 1.21

require (
	github.com/defioptimization/shared v0.0.0
	github.com/gin-contrib/cors v1.5.0
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/websocket v1.5.1
	github.com/stripe/stripe-go/v76 v76.0.0
)

replace github.com/defioptimization/shared => ../shared
