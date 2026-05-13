package resolver

import (
	"bufio"
	"context"
	"fmt"
	"net/netip"
	"os"
	"slices"
	"strings"

	"dns-manager/server/internal/domain"
)

type Resolver struct {
	path string
}

func New(path string) *Resolver {
	return &Resolver{path: path}
}

func (r *Resolver) Add(_ context.Context, address string) error {
	if _, err := netip.ParseAddr(address); err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInvalidAddr, address)
	}

	servers, err := r.parse()
	if err != nil {
		return err
	}

	if slices.Contains(servers, address) {
		return fmt.Errorf("%w: %s", domain.ErrAlreadyExists, address)
	}

	servers = append(servers, address)
	return r.write(servers)
}

func (r *Resolver) Remove(_ context.Context, address string) error {
	if _, err := netip.ParseAddr(address); err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInvalidAddr, address)
	}

	servers, err := r.parse()
	if err != nil {
		return err
	}

	idx := -1
	for i, s := range servers {
		if s == address {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("%w: %s", domain.ErrNotFound, address)
	}

	servers = append(servers[:idx], servers[idx+1:]...)
	return r.write(servers)
}

func (r *Resolver) List(_ context.Context) ([]string, error) {
	return r.parse()
}

func (r *Resolver) parse() ([]string, error) {
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open resolv.conf: %w", err)
	}
	defer f.Close()

	var servers []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "nameserver" {
			servers = append(servers, fields[1])
		}
	}

	return servers, scanner.Err()
}

func (r *Resolver) write(servers []string) error {
	var sb strings.Builder
	for _, s := range servers {
		sb.WriteString("nameserver ")
		sb.WriteString(s)
		sb.WriteByte('\n')
	}

	return os.WriteFile(r.path, []byte(sb.String()), 0644)
}
