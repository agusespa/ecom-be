package payload

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/agusespa/ecom-be/customer/internal/httperrors"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var contentType, errorMessage string
	var statusCode int

	if customErr, ok := err.(*httperrors.Error); ok {
		errorMessage = customErr.Message()
		statusCode = customErr.Status()
	} else {
		errorMessage = "Internal Server Error"
		statusCode = http.StatusInternalServerError
	}

	contentType = "text/plain"
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write([]byte(errorMessage)); err != nil {
		// TODO handle properly
		return
	}
}

func Write(w http.ResponseWriter, r *http.Request, payload any, cookies []*http.Cookie) {
	if payload == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	var contentType string
	var responseBytes []byte
	var err error

	switch p := payload.(type) {
	case string:
		contentType = "text/plain"
		responseBytes = []byte(p)
	case []byte:
		contentType = "application/octet-stream"
		responseBytes = p
	default:
		contentType = "application/json"
		responseBytes, err = json.Marshal(payload)
		if err != nil {
			WriteError(w, r, err)
			return
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(responseBytes)))

	if len(cookies) > 0 {
		for _, c := range cookies {
			http.SetCookie(w, c)
		}
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()
		gz := gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}
		w = gz
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(responseBytes); err != nil {
		WriteError(w, r, err)
		return
	}
}
