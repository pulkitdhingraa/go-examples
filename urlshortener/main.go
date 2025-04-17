package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var seedRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	shortUrlLength = 6
	charset = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type UrlShortener struct {
	shortToLong map[string]string
	mu sync.RWMutex
}

func NewUrlShortener() *UrlShortener {
	return &UrlShortener{
		shortToLong: make(map[string]string),
	}
}

func (us *UrlShortener) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post method allowed", http.StatusMethodNotAllowed)
		return
	}

	originalUrl := r.FormValue("url")
	if originalUrl == "" {
		http.Error(w, "url field is required", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()

	us.mu.Lock()
	us.shortToLong[shortKey] = originalUrl
	us.mu.Unlock()

	shortUrl := fmt.Sprintf("http://%s/short/%s", r.Host, shortKey)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(shortUrl))
	// fmt.Fprintf(w, shortUrl)
}

func (us *UrlShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/short/")
	us.mu.RLock()
	originalUrl, exists := us.shortToLong[key]
	us.mu.RUnlock()

	if !exists {
		http.NotFound(w,r)
		return
	}

	http.Redirect(w,r,originalUrl,http.StatusFound)
}

func generateShortKey() string {
	
	shortKey := make([]byte, shortUrlLength)
	for i := range shortKey {
		shortKey[i] = charset[seedRand.Intn(len(charset))]
	}
	return string(shortKey)
}

func main() {
	shortener := NewUrlShortener()

	http.HandleFunc("/shorten", shortener.ShortenUrl)
	http.HandleFunc("/short/", shortener.HandleRedirect)

	fmt.Println("Server started on localhost 8080")
	http.ListenAndServe(":8080", nil)
}
