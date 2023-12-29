package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	var err error

	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s |%s %d %s| %s |%s %s %s %s | %s | %s | %s\n",
				param.TimeStamp.Format(time.RFC1123),
				param.StatusCodeColor(),
				param.StatusCode,
				param.ResetColor(),
				param.ClientIP,
				param.MethodColor(),
				param.Method,
				param.ResetColor(),
				param.Path,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		}),
	)

	router.Static("/static", "./static")
	router.Any("/", Wrap(http.FileServer(http.Dir("./"))))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	var (
		listener net.Listener
	)
	if listener, err = net.Listen("tcp", srv.Addr); err != nil {
		logrus.Fatalf("listen: %s\n", err)
	}

	logrus.Infof("Starting Server! Listening on: %s\n", srv.Addr)

	srv.Serve(listener)

}

func Wrap(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
