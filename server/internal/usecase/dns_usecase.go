package usecase

import (
	"context"
	"fmt"
	"net/netip"

	"dns-manager/server/internal/domain"
)

type DnsUseCase struct {
	repo domain.DnsRepository
}

func NewDnsUseCase(repo domain.DnsRepository) *DnsUseCase {
	return &DnsUseCase{repo: repo}
}

func (uc *DnsUseCase) Add(ctx context.Context, address string) error {
	if _, err := netip.ParseAddr(address); err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInvalidAddr, address)
	}
	return uc.repo.Add(ctx, address)
}

func (uc *DnsUseCase) Remove(ctx context.Context, address string) error {
	if _, err := netip.ParseAddr(address); err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInvalidAddr, address)
	}
	return uc.repo.Remove(ctx, address)
}

func (uc *DnsUseCase) List(ctx context.Context) ([]string, error) {
	return uc.repo.List(ctx)
}
