package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = io.WriteString(w, "Hello astaxie!\n")
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var tokenBytes [32]byte
		if _, err := rand.Read(tokenBytes[:]); err != nil {
			http.Error(w, "could not prepare login form", http.StatusInternalServerError)
			return
		}
		token := base64.RawURLEncoding.EncodeToString(tokenBytes[:])
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf",
			Value:    token,
			Path:     "/login",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   300,
		})
		w.Header().Set("Cache-Control", "no-store")
		t, err := template.ParseFiles("login.html")
		if err != nil {
			http.Error(w, "login form unavailable", http.StatusInternalServerError)
			return
		}
		if err := t.Execute(w, token); err != nil {
			log.Printf("render login form: %v", err)
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid login request", http.StatusBadRequest)
			return
		}
		cookie, err := r.Cookie("csrf")
		if err != nil || subtle.ConstantTimeCompare([]byte(r.Form.Get("token")), []byte(cookie.Value)) != 1 {
			http.Error(w, "invalid login request", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = io.WriteString(w, "Authentication is not implemented in this sample.\n")
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
