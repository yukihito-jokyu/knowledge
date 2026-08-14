//go:build !integrationtest

package sqlite

import "context"

func waitIntegrationGate(context.Context, string) error {
	return nil
}
