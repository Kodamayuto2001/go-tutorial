package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginDoesNotReflectSubmittedCredentials(t *testing.T) {
	const token = "test-csrf-token"
	form := url.Values{
		"token":    {token},
		"username": {"<script>alert(1)</script>"},
		"password": {"sample-only-password"},
	}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "csrf", Value: token})
	rec := httptest.NewRecorder()

	login(rec, req)

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotImplemented)
	}
	if strings.Contains(rec.Body.String(), form.Get("username")) || strings.Contains(rec.Body.String(), form.Get("password")) {
		t.Fatal("login response reflected submitted values")
	}
}

func TestLoginRejectsMismatchedCSRFToken(t *testing.T) {
	form := url.Values{"token": {"submitted"}, "username": {"user"}, "password": {"sample-only-password"}}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "csrf", Value: "different"})
	rec := httptest.NewRecorder()

	login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
