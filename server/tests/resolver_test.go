package tests

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dns-manager/server/internal/adapter/resolver"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "resolv.conf")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAdd(t *testing.T) {
	path := writeTempFile(t, "nameserver 1.1.1.1\n")
	r := resolver.New(path)
	ctx := context.Background()

	if err := r.Add(ctx, "8.8.8.8"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	servers, err := r.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 2 || servers[1] != "8.8.8.8" {
		t.Fatalf("expected [1.1.1.1 8.8.8.8], got %v", servers)
	}
}

func TestAdd_Duplicate(t *testing.T) {
	path := writeTempFile(t, "nameserver 8.8.8.8\n")
	r := resolver.New(path)

	if err := r.Add(context.Background(), "8.8.8.8"); err == nil {
		t.Fatal("expected error for duplicate, got nil")
	}
}

func TestAdd_InvalidIP(t *testing.T) {
	path := writeTempFile(t, "")
	r := resolver.New(path)

	if err := r.Add(context.Background(), "not-an-ip"); err == nil {
		t.Fatal("expected error for invalid IP, got nil")
	}
}

func TestRemove(t *testing.T) {
	path := writeTempFile(t, "nameserver 1.1.1.1\nnameserver 8.8.8.8\n")
	r := resolver.New(path)
	ctx := context.Background()

	if err := r.Remove(ctx, "1.1.1.1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	servers, err := r.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 || servers[0] != "8.8.8.8" {
		t.Fatalf("expected [8.8.8.8], got %v", servers)
	}
}

func TestRemove_NonExistent(t *testing.T) {
	path := writeTempFile(t, "nameserver 8.8.8.8\n")
	r := resolver.New(path)

	if err := r.Remove(context.Background(), "1.1.1.1"); err == nil {
		t.Fatal("expected error for non-existent, got nil")
	}
}

func TestRemove_InvalidIP(t *testing.T) {
	path := writeTempFile(t, "")
	r := resolver.New(path)

	if err := r.Remove(context.Background(), "bad"); err == nil {
		t.Fatal("expected error for invalid IP, got nil")
	}
}

func TestList(t *testing.T) {
	path := writeTempFile(t, "nameserver 8.8.8.8\nnameserver 1.1.1.1\n")
	r := resolver.New(path)

	servers, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("expected 2 servers, got %d: %v", len(servers), servers)
	}
	if servers[0] != "8.8.8.8" || servers[1] != "1.1.1.1" {
		t.Fatalf("expected [8.8.8.8 1.1.1.1], got %v", servers)
	}
}

func TestList_Empty(t *testing.T) {
	path := writeTempFile(t, "")
	r := resolver.New(path)

	servers, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(servers) != 0 {
		t.Fatalf("expected empty list, got %v", servers)
	}
}

func TestList_NonExistentFile(t *testing.T) {
	r := resolver.New("/tmp/nonexistent-resolv-test")

	servers, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(servers) != 0 {
		t.Fatalf("expected empty list, got %v", servers)
	}
}

func TestParseSkipsComments(t *testing.T) {
	content := "# this is a comment\nnameserver 8.8.8.8\n# another comment\n"
	path := writeTempFile(t, content)
	r := resolver.New(path)

	servers, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(servers) != 1 || servers[0] != "8.8.8.8" {
		t.Fatalf("expected [8.8.8.8], got %v", servers)
	}
}

func TestParseSkipsEmptyLines(t *testing.T) {
	content := "\n\nnameserver 8.8.8.8\n\n"
	path := writeTempFile(t, content)
	r := resolver.New(path)

	servers, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(servers) != 1 || servers[0] != "8.8.8.8" {
		t.Fatalf("expected [8.8.8.8], got %v", servers)
	}
}

func TestPreservesNonNameserverLines(t *testing.T) {
	content := "search example.com\nnameserver 1.1.1.1\noptions timeout:2\n"
	path := writeTempFile(t, content)
	r := resolver.New(path)

	if err := r.Add(context.Background(), "8.8.8.8"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	result := string(data)
	if !strings.Contains(result, "search example.com") {
		t.Fatal("search directive was lost")
	}
	if !strings.Contains(result, "options timeout:2") {
		t.Fatal("options directive was lost")
	}
}

func TestFullWorkflow(t *testing.T) {
	path := writeTempFile(t, "")
	r := resolver.New(path)
	ctx := context.Background()

	servers, err := r.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 0 {
		t.Fatalf("expected empty, got %v", servers)
	}

	if err := r.Add(ctx, "8.8.8.8"); err != nil {
		t.Fatal(err)
	}
	if err := r.Add(ctx, "1.1.1.1"); err != nil {
		t.Fatal(err)
	}

	servers, err = r.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 2 {
		t.Fatalf("expected 2, got %v", servers)
	}

	if err := r.Remove(ctx, "8.8.8.8"); err != nil {
		t.Fatal(err)
	}
	servers, err = r.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 || servers[0] != "1.1.1.1" {
		t.Fatalf("expected [1.1.1.1], got %v", servers)
	}
}
