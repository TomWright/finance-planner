package domain_test

import (
	"github.com/tomwright/finance-planner/internal/application/domain"
	"testing"
)

func TestTransactionCollection_Sum(t *testing.T) {
	t.Parallel()

	c := domain.NewTransactionCollection()

	c.Add(
		domain.NewTransaction().WithAmount(100),
		domain.NewTransaction().WithAmount(1),
		domain.NewTransaction().WithAmount(-5),
		domain.NewTransaction().WithAmount(201),
		domain.NewTransaction().WithAmount(-200),
	)

	sum := c.Sum()
	if exp, got := int64(97), sum; exp != got {
		t.Errorf("expected sum %d, got %d", exp, got)
	}
}
