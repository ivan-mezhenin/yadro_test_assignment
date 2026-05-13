package tests

import (
	"context"
	"errors"
	"testing"

	"dns-manager/server/internal/domain"
	"dns-manager/server/internal/usecase"
)

type mockRepo struct {
	servers []string
	addErr  error
	remErr  error
	listErr error
}

func (m *mockRepo) Add(_ context.Context, address string) error {
	if m.addErr != nil {
		return m.addErr
	}
	for _, s := range m.servers {
		if s == address {
			return domain.ErrAlreadyExists
		}
	}
	m.servers = append(m.servers, address)
	return nil
}

func (m *mockRepo) Remove(_ context.Context, address string) error {
	if m.remErr != nil {
		return m.remErr
	}
	for i, s := range m.servers {
		if s == address {
			m.servers = append(m.servers[:i], m.servers[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockRepo) List(_ context.Context) ([]string, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.servers, nil
}

func TestUseCase_Add_Valid(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{})
	err := uc.Add(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestUseCase_Add_InvalidIP(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{})
	err := uc.Add(context.Background(), "not-an-ip")
	if !errors.Is(err, domain.ErrInvalidAddr) {
		t.Fatal("expected ErrInvalidAddr")
	}
}

func TestUseCase_Add_Duplicate(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{servers: []string{"8.8.8.8"}})
	err := uc.Add(context.Background(), "8.8.8.8")
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatal("expected ErrAlreadyExists")
	}
}

func TestUseCase_Remove_Valid(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{servers: []string{"8.8.8.8"}})
	err := uc.Remove(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestUseCase_Remove_InvalidIP(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{})
	err := uc.Remove(context.Background(), "bad")
	if !errors.Is(err, domain.ErrInvalidAddr) {
		t.Fatal("expected ErrInvalidAddr")
	}
}

func TestUseCase_Remove_NonExistent(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{servers: []string{"1.1.1.1"}})
	err := uc.Remove(context.Background(), "8.8.8.8")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatal("expected ErrNotFound")
	}
}

func TestUseCase_List(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{servers: []string{"8.8.8.8", "1.1.1.1"}})
	servers, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("expected 2, got %d", len(servers))
	}
}

func TestUseCase_List_Empty(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{})
	servers, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(servers) != 0 {
		t.Fatalf("expected 0, got %d", len(servers))
	}
}

func TestUseCase_List_RepoError(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{listErr: errors.New("disk error")})
	servers, err := uc.List(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if servers != nil {
		t.Fatal("expected nil servers on error")
	}
}

func TestUseCase_Add_RepoError(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{addErr: errors.New("permission denied")})
	err := uc.Add(context.Background(), "8.8.8.8")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, domain.ErrAlreadyExists) || errors.Is(err, domain.ErrInvalidAddr) {
		t.Fatal("expected raw repo error, not domain error")
	}
}

func TestUseCase_Remove_RepoError(t *testing.T) {
	uc := usecase.NewDnsUseCase(&mockRepo{remErr: errors.New("permission denied")})
	err := uc.Remove(context.Background(), "8.8.8.8")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrInvalidAddr) {
		t.Fatal("expected raw repo error, not domain error")
	}
}
