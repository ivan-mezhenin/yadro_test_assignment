package resolver

import (
	"bufio"
	"context"
	"fmt"
	"net/netip"
	"os"
	"strings"
	"sync"

	"dns-manager/server/internal/domain"
)

type parsedFile struct {
	lines   []string
	servers []string
}

type Resolver struct {
	path string
	mu   *sync.Mutex
}

func New(path string) *Resolver {
	return &Resolver{path: path, mu: &sync.Mutex{}}
}

func (r *Resolver) Add(_ context.Context, address string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := netip.ParseAddr(address); err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInvalidAddr, address)
	}

	pf, err := r.parse()
	if err != nil {
		return err
	}

	for _, s := range pf.servers {
		if s == address {
			return fmt.Errorf("%w: %s", domain.ErrAlreadyExists, address)
		}
	}

	pf.servers = append(pf.servers, address)
	return r.write(pf)
}

func (r *Resolver) Remove(_ context.Context, address string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := netip.ParseAddr(address); err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInvalidAddr, address)
	}

	pf, err := r.parse()
	if err != nil {
		return err
	}

	idx := -1
	for i, s := range pf.servers {
		if s == address {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("%w: %s", domain.ErrNotFound, address)
	}

	pf.servers = append(pf.servers[:idx], pf.servers[idx+1:]...)
	return r.write(pf)
}

func (r *Resolver) List(_ context.Context) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pf, err := r.parse()
	if err != nil {
		return nil, err
	}
	return pf.servers, nil
}

func (r *Resolver) parse() (*parsedFile, error) {
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &parsedFile{}, nil
		}
		return nil, fmt.Errorf("open resolv.conf: %w", err)
	}
	defer f.Close()

	var pf parsedFile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		pf.lines = append(pf.lines, line)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) == 2 && fields[0] == "nameserver" {
			pf.servers = append(pf.servers, fields[1])
		}
	}
	return &pf, scanner.Err()
}

func (r *Resolver) write(pf *parsedFile) error {
	var sb strings.Builder
	si := 0
	for _, line := range pf.lines {
		trimmed := strings.TrimSpace(line)
		fields := strings.Fields(trimmed)
		if len(fields) == 2 && fields[0] == "nameserver" {
			if si < len(pf.servers) {
				sb.WriteString("nameserver ")
				sb.WriteString(pf.servers[si])
				sb.WriteByte('\n')
				si++
			}
		} else {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	}
	for ; si < len(pf.servers); si++ {
		sb.WriteString("nameserver ")
		sb.WriteString(pf.servers[si])
		sb.WriteByte('\n')
	}

	return os.WriteFile(r.path, []byte(sb.String()), 0644)
}
