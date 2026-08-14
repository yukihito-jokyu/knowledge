package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yukihito-jokyu/knowledge/internal/persistence/sqlite"
)

var (
	userConfigDir      = os.UserConfigDir
	makeStoreDirectory = os.MkdirAll
	openSQLiteStore    = sqlite.Open
)

// executeRetrieval はDEC-FEAT-005の既定Storeを開いて取得操作を実行する。
func executeRetrieval(ctx context.Context, parsed command) (any, cliError, bool) {
	store, err := openDefaultRetrievalStore(ctx)
	if err != nil {
		return nil, cliError{
			code:    storageError,
			message: "Knowledge Storeを開けません",
		}, true
	}
	defer store.Close()

	return executeRetrievalWithStore(ctx, parsed, store)
}

func openDefaultRetrievalStore(ctx context.Context) (*sqlite.Store, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	path, err := defaultStorePath()
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := makeStoreDirectory(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create knowledge store directory: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return openSQLiteStore(ctx, path)
}

func defaultStorePath() (string, error) {
	directory, err := userConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user configuration directory: %w", err)
	}

	return filepath.Join(directory, "knowledge", "knowledge.db"), nil
}
