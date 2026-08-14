package main

import (
	"context"
	"errors"
	"testing"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

type createStoreStub struct {
	request domain.CreateRequest
	err     error
}

func (s *createStoreStub) CreateAssertion(ctx context.Context, request domain.CreateRequest) (domain.CreateResult, error) {
	s.request = request
	if s.err != nil {
		return domain.CreateResult{}, s.err
	}

	return domain.CreateResult{
		AssertionID: "asrt_01",
		Revision:    1,
		EvidenceIDs: []string{"evd_01"},
		Concepts: []domain.Concept{
			{
				ID:   "cpt_01",
				Name: "channel",
			},
		},
		RelationIDs: []string{"rel_01"},
	}, nil
}

func TestExecuteCreateWithStore(t *testing.T) {
	parsed, err := parseCommand([]string{
		"create",
		"--normalized-text", "channel send",
		"--scope-key", "language",
		"--scope-value", "Go",
		"--concept", "channel",
		"--concept-alias", "chan",
		"--alias-kind", "identifier",
		"--alias-value", "ch",
		"--relation-type", "causes",
		"--relation-target-kind", "assertion",
		"--relation-target-id", "asrt_00",
		"--evidence-kind", "user_code",
		"--evidence-text", "send(ch)",
		"--evidence-observed-at", "2026-08-14T00:00:00Z",
		"--valid-from", "2026-08-14T00:00:00Z",
		"--version-scope", "1.26",
	})
	if err.code != "" {
		t.Fatalf("createをparseする: %v", err)
	}
	store := &createStoreStub{}
	data, executionError, handled := executeCreateWithStore(context.Background(), parsed, store)
	if !handled || executionError.code != "" || data == nil {
		t.Fatalf("executeCreateWithStore() = %#v, %#v, %t", data, executionError, handled)
	}
	if store.request.Temporal == nil || store.request.Temporal.ValidFrom == nil || *store.request.Temporal.ValidFrom != "2026-08-14T00:00:00.000000000Z" || len(store.request.Concepts[0].Aliases) != 1 {
		t.Fatalf("create request = %#v", store.request)
	}
	for _, createErr := range []struct {
		err  error
		code errorCode
	}{
		{
			err:  domain.ErrCreateConflict,
			code: conflictError,
		},
		{
			err:  domain.ErrCreateRelationTargetNotFound,
			code: notFoundError,
		},
		{
			err:  errors.New("storage failure"),
			code: storageError,
		},
	} {
		store.err = createErr.err
		_, executionError, handled = executeCreateWithStore(context.Background(), parsed, store)
		if !handled || executionError.code != createErr.code {
			t.Fatalf("create error = %#v, handled=%t", executionError, handled)
		}
	}
	if _, _, handled := executeCreateWithStore(context.Background(), command{
		operation: "unknown",
	}, store); handled {
		t.Fatal("未知操作をcreateとして扱いました")
	}
}

func TestExecuteCommandCreate(t *testing.T) {
	originalUserConfigDir := userConfigDir
	t.Cleanup(func() {
		userConfigDir = originalUserConfigDir
	})
	parsed, err := parseCommand([]string{
		"create",
		"--normalized-text", "text",
		"--evidence-kind", "user_code",
		"--evidence-text", "evidence",
		"--evidence-observed-at", "2026-08-14T00:00:00Z",
	})
	if err.code != "" {
		t.Fatalf("createをparseする: %v", err)
	}
	userConfigDir = func() (string, error) {
		return t.TempDir(), nil
	}
	if data, executionError, handled := executeCommand(context.Background(), parsed); !handled || executionError.code != "" || data == nil {
		t.Fatalf("executeCommand() = %#v, %#v, %t", data, executionError, handled)
	}
	userConfigDir = func() (string, error) {
		return "", errors.New("no config directory")
	}
	if _, executionError, handled := executeCommand(context.Background(), parsed); !handled || executionError.code != storageError {
		t.Fatalf("Store open failure = %#v, handled=%t", executionError, handled)
	}
	userConfigDir = func() (string, error) {
		return t.TempDir(), nil
	}
	search, searchErr := parseCommand([]string{
		"search-text",
		"--query", "text",
	})
	if searchErr.code != "" {
		t.Fatalf("searchをparseする: %v", searchErr)
	}
	if _, executionError, handled := executeCommand(context.Background(), search); !handled || executionError.code != "" {
		t.Fatalf("retrieval delegation = %#v, handled=%t", executionError, handled)
	}
}

func TestExecuteCommandRejectsUnregisteredParsedOperation(t *testing.T) {
	parsed, err := parseCommand([]string{
		"attach-evidence",
		"--assertion-id", "asrt_01",
		"--evidence-kind", "user_code",
		"--evidence-text", "evidence",
		"--evidence-observed-at", "2026-08-14T00:00:00Z",
	})
	if err.code != "" {
		t.Fatalf("attach-evidenceをparseする: %v", err)
	}

	if _, executionError, handled := executeCommand(context.Background(), parsed); handled || executionError.code != "" {
		t.Fatalf("executeCommand() = %#v, handled=%t", executionError, handled)
	}
}

func TestCreateRequestWithoutTemporal(t *testing.T) {
	parsed, err := parseCommand([]string{
		"create",
		"--normalized-text", "text",
		"--evidence-kind", "user_code",
		"--evidence-text", "evidence",
		"--evidence-observed-at", "2026-08-14T00:00:00Z",
	})
	if err.code != "" {
		t.Fatalf("createをparseする: %v", err)
	}
	request, requestErr := createRequest(parsed.options)
	if requestErr.code != "" || request.Temporal != nil || len(request.Aliases) != 0 || len(request.Relations) != 0 {
		t.Fatalf("createRequest() = %#v, %#v", request, requestErr)
	}
	invalid := command{
		operation: "create",
		options: []option{
			{
				name:  "normalized-text",
				value: "text",
			},
			{
				name:  "valid-from",
				value: "invalid",
			},
		},
	}
	if _, executionError, handled := executeCreateWithStore(context.Background(), invalid, &createStoreStub{}); !handled || executionError.code != validationError {
		t.Fatalf("invalid temporal = %#v, handled=%t", executionError, handled)
	}
}

func TestCreateRejectsOrphanedAliasAndRelationOptions(t *testing.T) {
	base := []string{
		"create",
		"--normalized-text", "text",
		"--evidence-kind", "user_code",
		"--evidence-text", "evidence",
		"--evidence-observed-at", "2026-08-14T00:00:00Z",
	}
	for _, options := range [][]string{
		{"--alias-value", "name"},
		{"--relation-target-kind", "assertion"},
		{"--relation-target-id", "asrt_01"},
	} {
		arguments := append(append([]string{}, base...), options...)
		_, err := parseCommand(arguments)
		if err.code != validationError {
			t.Fatalf("孤立option error = %#v", err)
		}
	}
}

func TestCreateRejectsInvalidConceptAliasGrouping(t *testing.T) {
	for _, arguments := range [][]string{
		{"create", "--normalized-text", "text", "--concept", "foo", "--concept-alias", "foo", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"},
		{"create", "--normalized-text", "text", "--concept", "foo", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z", "--concept-alias", "bar"},
	} {
		_, err := parseCommand(arguments)
		if err.code != validationError {
			t.Fatalf("concept alias error = %#v", err)
		}
	}
	if hasPreviousConcept([]option{{name: "concept-alias"}}, 0) {
		t.Fatal("先行Conceptなしを許可しました")
	}
}
