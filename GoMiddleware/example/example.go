// GoMiddleware project example.go
package main

import (
	"bytes"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/GrassInWind2019/GoDemo/GoMiddleware/httpLog"
	"github.com/goji/httpauth"
)

func enforceXMLHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		if http.DetectContentType(buf.Bytes()) != "text/xml; charset=utf-8" {
			http.Error(w, http.StatusText(415), 415)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logHandler(h http.Handler) http.Handler {
	logFile, err := os.OpenFile("httpServerLog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return httpLog.HttpLogHandler(logFile, h)
}

func main() {
	authHandler := httpauth.SimpleBasicAuth("username", "password")

	http.Handle("/", logHandler(authHandler(enforceXMLHandler(http.HandlerFunc(rootReqHandler)))))
	http.ListenAndServe(":10000", nil)
}

func rootReqHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello! GrassInWind2019!"))
	rand.Seed(time.Now().UnixNano())
	workTime := rand.Intn(100)
	//simulate working
	time.Sleep(time.Duration(time.Duration(workTime) * time.Millisecond))
}
