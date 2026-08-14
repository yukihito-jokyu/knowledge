//go:build integrationtest

package sqlite

import (
	"context"
	"fmt"
	"os"
)

const (
	integrationGateStageEnvironment = "KNOWLEDGE_TEST_INTEGRATION_GATE_STAGE"
	integrationGateReadyEnvironment = "KNOWLEDGE_TEST_INTEGRATION_GATE_READY"
)

// waitIntegrationGateはprocess境界testだけで処理中の割込みを同期する。
func waitIntegrationGate(ctx context.Context, stage string) error {
	if os.Getenv(integrationGateStageEnvironment) != stage {
		return nil
	}
	path := os.Getenv(integrationGateReadyEnvironment)
	if path == "" {
		return fmt.Errorf("integration gateのreadyファイルがありません")
	}
	if err := os.WriteFile(path, []byte(stage), 0o600); err != nil {
		return fmt.Errorf("integration gateのreadyファイル作成: %w", err)
	}
	<-ctx.Done()

	return ctx.Err()
}
