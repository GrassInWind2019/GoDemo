# GoMiddleware基础知识  
编写Golang net/http库的中间件需要满足http.Handler这个接口  
```
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
可以使用http.HandlerFunc对自己实现的函数进行转换以满足http.Handler接口, 从源码可见， httpHandlerFunc实现了ServeHTTP方法，并调用HandlerFunc,  
也就是调用用户实现的函数。  
```
// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```
可以通过如下的结构，来实现了一个中间件handler链，以实现调用多个中间件  
```
func exampleMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Our middleware logic goes here...
    next.ServeHTTP(w, r)
  })
}
```
## Go log  
本文实现了一个日志中间件Demo， 主要是实现一个日志中间件handler和一个http.ResponseWriter的包装接口，从而获取处理结果并记录到日志中    
日志中间件handler如下  
```
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
```
http.ResponseWriter的包装如下  
```
type logResponseWriter struct {
	W      http.ResponseWriter
	Status *int
	Size   *int
}
func (l logResponseWriter) Write(b []byte) (int, error) {
	size, err := l.W.Write(b)
  //记录response大小
	(*l.Size) += size
	return size, err
}

func (l logResponseWriter) WriteHeader(s int) {
	l.W.WriteHeader(s)
  //记录response状态,以便后面记入日志中
	(*l.Status) = s
}
```
## Demo 结果
### log  
```
http request version=HTTP/1.1, method=POST, host=localhost:10000, uri=/, receive timestamp=Sun, 11 Aug 2019 17:43:33 CST, handle cost time=73.0041ms, status=200, size=23
http request version=HTTP/1.1, method=POST, host=localhost:10000, uri=/, receive timestamp=Sun, 11 Aug 2019 17:44:19 CST, handle cost time=0s, status=401, size=13
http request version=HTTP/1.1, method=GET, host=localhost:10000, uri=/, receive timestamp=Sun, 11 Aug 2019 17:45:04 CST, handle cost time=0s, status=400, size=12
http request version=HTTP/1.1, method=POST, host=localhost:10000, uri=/, receive timestamp=Sun, 11 Aug 2019 17:45:34 CST, handle cost time=0s, status=401, size=13
```
### curl 测试截图    
[curl_test1.png](https://github.com/GrassInWind2019/GoDemo/blob/master/GoMiddleware/curl_test1.png)
[curl_test2.png](https://github.com/GrassInWind2019/GoDemo/blob/master/GoMiddleware/curl_test2.png)
