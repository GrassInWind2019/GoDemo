package httpLog

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type LogFormatParams struct {
	Req       *http.Request
	StartTime time.Time
	Cost      time.Duration
	status    int
	size      int
}

type logResponseWriter struct {
	W      http.ResponseWriter
	Status *int
	Size   *int
}

func (l logResponseWriter) Header() http.Header {
	return l.W.Header()
}

func (l logResponseWriter) Write(b []byte) (int, error) {
	size, err := l.W.Write(b)
	(*l.Size) += size
	return size, err
}

func (l logResponseWriter) WriteHeader(s int) {
	l.W.WriteHeader(s)
	(*l.Status) = s
}

func HttpLogHandler(out io.Writer, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		status := http.StatusOK
		size := 0
		logResWriter := logResponseWriter{
			W:      w,
			Status: &status,
			Size:   &size,
		}
		h.ServeHTTP(logResWriter, r)
		logParams := LogFormatParams{
			Req:       r,
			StartTime: t,
			Cost:      time.Duration(time.Since(t).Nanoseconds()),
			status:    (*logResWriter.Status),
			size:      (*logResWriter.Size),
		}
		fmt.Println(logParams)
		log := logFormat(&logParams)
		_, err := out.Write(log)
		if err != nil {
			fmt.Println("Write log failed")
		}
	}

	return http.HandlerFunc(fn)
}

func logFormat(logParams *LogFormatParams) []byte {
	log := fmt.Sprintf("http request version=%s, method=%s, host=%s, uri=%s, receive timestamp=%s, handle cost time=%s, status=%d, size=%d\n",
		logParams.Req.Proto,
		logParams.Req.Method,
		logParams.Req.Host,
		logParams.Req.URL.String(),
		logParams.StartTime.Format(time.RFC1123),
		logParams.Cost.String(),
		logParams.status,
		logParams.size)
	fmt.Println(log)
	var logSlice []byte
	logSlice = append(logSlice, log...)
	//fmt.Println("logFormat: ", logSlice)
	return logSlice
}
