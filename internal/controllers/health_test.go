package controllers

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type fakeDatabase struct{ err error }

func (f fakeDatabase) Ping() error { return f.err }

func checkHealth(t *testing.T, database Pinger) (int, string) {
	t.Helper()
	app := fiber.New()
	controller := NewHealthController(database)
	app.Get("/health", controller.Check)

	response, err := app.Test(httptest.NewRequest("GET", "/health", nil))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(response.Body)
	return response.StatusCode, string(body)
}

func TestHealthIsOkWhenTheDatabaseAnswers(t *testing.T) {
	status, body := checkHealth(t, fakeDatabase{})
	if status != 200 || body != `{"status":"ok"}` {
		t.Fatalf("got %d %s", status, body)
	}
}

func TestHealthIsUnavailableWhenTheDatabaseFails(t *testing.T) {
	status, body := checkHealth(t, fakeDatabase{err: errors.New("table unreachable")})
	if status != 503 || body != `{"status":"unavailable"}` {
		t.Fatalf("got %d %s", status, body)
	}
}
