package main

import (
	"log"
	"net/http"
	"time"
)

// proxies the ResponseWriter and keeps track of the status and bytes written
type responseObserver struct {
	http.ResponseWriter
	// HTTP status of this response
	status int
	// number of bytes written
	written int64
	// internal field to keep track if status header has been written yet
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func requestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := &responseObserver{
			ResponseWriter: w,
		}

		targetMux.ServeHTTP(writer, r)

		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t%d\t\t%s\t\t%s\t\t%v",
			r.Method,
			writer.status,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
	})
}
