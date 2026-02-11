package integration

import (
	"context"
	"testing"
)

func TestIntegration_ListBlockchains(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	blockchains, err := client.Blockchains().ListBlockchains(ctx, nil)
	if err != nil {
		t.Fatalf("ListBlockchains() error = %v", err)
	}

	t.Logf("Found %d blockchains", len(blockchains))
	for _, b := range blockchains {
		t.Logf("Blockchain: Symbol=%s, Name=%s, Network=%s", b.Symbol, b.Name, b.Network)
	}
}
