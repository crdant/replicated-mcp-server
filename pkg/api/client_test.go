package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crdant/replicated-mcp-server/pkg/config"
	"github.com/crdant/replicated-mcp-server/pkg/logging"
)

const (
	testAppID = "app1"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
	}
	logger := logging.NewLogger("error")

	client, err := NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
}

func TestNewClientWithInvalidConfig(t *testing.T) {
	cfg := &config.Config{
		APIToken: "",
		Timeout:  30 * time.Second,
	}
	logger := logging.NewLogger("error")

	_, err := NewClient(cfg, logger)
	if err == nil {
		t.Fatal("Expected error for empty API token")
	}
}

func TestListApplications(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/apps" {
			t.Errorf("Expected path /v3/apps, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"apps": [{"id": "app1", "name": "Test App", "slug": "test-app"}]}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
		Endpoint: server.URL,
	}
	logger := logging.NewLogger("error")

	client, err := NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	apps, err := client.ListApplications(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("Failed to list applications: %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("Expected 1 app, got %d", len(apps))
	}
	if apps[0].ID != testAppID {
		t.Errorf("Expected app ID 'app1', got %s", apps[0].ID)
	}
}

func TestGetApplication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/app/test-app" {
			t.Errorf("Expected path /v3/app/test-app, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"app": {"id": "app1", "name": "Test App", "slug": "test-app"}}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
		Endpoint: server.URL,
	}
	logger := logging.NewLogger("error")

	client, err := NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	app, err := client.GetApplication(context.Background(), "test-app")
	if err != nil {
		t.Fatalf("Failed to get application: %v", err)
	}

	if app.ID != testAppID {
		t.Errorf("Expected app ID 'app1', got %s", app.ID)
	}
	if app.Slug != "test-app" {
		t.Errorf("Expected app slug 'test-app', got %s", app.Slug)
	}
}

func TestListReleases(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/app/test-app/releases" {
			t.Errorf("Expected path /v3/app/test-app/releases, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"releases": [{"sequence": 1, "version": "1.0.0"}]}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
		Endpoint: server.URL,
	}
	logger := logging.NewLogger("error")

	client, err := NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	releases, err := client.ListReleases(context.Background(), "test-app", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list releases: %v", err)
	}

	if len(releases) != 1 {
		t.Fatalf("Expected 1 release, got %d", len(releases))
	}
	if releases[0].Sequence != 1 {
		t.Errorf("Expected release sequence 1, got %d", releases[0].Sequence)
	}
}

func TestListChannels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/app/test-app/channels" {
			t.Errorf("Expected path /v3/app/test-app/channels, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"channels": [{"id": "ch1", "name": "Stable", "slug": "stable"}]}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
		Endpoint: server.URL,
	}
	logger := logging.NewLogger("error")

	client, err := NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	channels, err := client.ListChannels(context.Background(), "test-app", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list channels: %v", err)
	}

	if len(channels) != 1 {
		t.Fatalf("Expected 1 channel, got %d", len(channels))
	}
	if channels[0].ID != "ch1" {
		t.Errorf("Expected channel ID 'ch1', got %s", channels[0].ID)
	}
}

func TestListCustomers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/app/test-app/customers" {
			t.Errorf("Expected path /v3/app/test-app/customers, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"customers": [{"id": "cust1", "name": "Test Customer"}]}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIToken: "test-token",
		Timeout:  30 * time.Second,
		Endpoint: server.URL,
	}
	logger := logging.NewLogger("error")

	client, err := NewClient(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	customers, err := client.ListCustomers(context.Background(), "test-app", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list customers: %v", err)
	}

	if len(customers) != 1 {
		t.Fatalf("Expected 1 customer, got %d", len(customers))
	}
	if customers[0].ID != "cust1" {
		t.Errorf("Expected customer ID 'cust1', got %s", customers[0].ID)
	}
}
