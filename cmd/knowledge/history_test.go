package main

import (
	"context"
	"errors"
	"testing"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

type historyStoreStub struct{ err error }

func (s historyStoreStub) AttachEvidence(context.Context, domain.AttachEvidenceRequest) (domain.AttachEvidenceResult, error) {
	if s.err != nil {
		return domain.AttachEvidenceResult{}, s.err
	}

	return domain.AttachEvidenceResult{
		AssertionID: "asrt_01",
		EvidenceID:  "evd_02",
	}, nil
}

func (s historyStoreStub) ReviseAssertion(context.Context, domain.ReviseRequest) (domain.ReviseResult, error) {
	if s.err != nil {
		return domain.ReviseResult{}, s.err
	}

	return domain.ReviseResult{
		AssertionID:      "asrt_01",
		PreviousRevision: 1,
		Revision:         2,
	}, nil
}

func (s historyStoreStub) Supersede(context.Context, domain.SupersedeRequest) (domain.SupersedeResult, error) {
	if s.err != nil {
		return domain.SupersedeResult{}, s.err
	}

	return domain.SupersedeResult{
		RelationID:             "rel_01",
		SupersededAssertionID:  "asrt_01",
		ReplacementAssertionID: "asrt_02",
	}, nil
}

func TestExecuteHistoryWithStore(t *testing.T) {
	for _, arguments := range [][]string{
		{"attach-evidence", "--assertion-id", "asrt_01", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"},
		{"revise", "--assertion-id", "asrt_01", "--normalized-text", "revised"},
		{"supersede", "--superseded-assertion-id", "asrt_01", "--replacement-assertion-id", "asrt_02"},
	} {
		parsed, err := parseCommand(arguments)
		if err.code != "" {
			t.Fatalf("parseCommand(%v) = %v", arguments, err)
		}
		if _, cliErr, handled := executeHistoryWithStore(context.Background(), parsed, historyStoreStub{}); !handled || cliErr.code != "" {
			t.Fatalf("execute = %#v, %t", cliErr, handled)
		}
		for _, storeErr := range []error{domain.ErrAssertionNotFound, domain.ErrMutationConflict, errors.New("storage")} {
			if _, cliErr, handled := executeHistoryWithStore(context.Background(), parsed, historyStoreStub{err: storeErr}); !handled || cliErr.code == "" {
				t.Fatalf("error=%v: %#v, %t", storeErr, cliErr, handled)
			}
		}
	}
	if _, _, handled := executeHistoryWithStore(context.Background(), command{operation: "unknown"}, historyStoreStub{}); handled {
		t.Fatal("未知操作を履歴更新として扱いました")
	}
	if _, cliErr, handled := executeHistoryWithStore(context.Background(), command{
		operation: "revise",
		options: []option{
			{
				name:  "assertion-id",
				value: "asrt_01",
			},
			{
				name:  "normalized-text",
				value: "revised",
			},
			{
				name:  "observed-at",
				value: "invalid",
			},
		},
	}, historyStoreStub{}); !handled || cliErr.code != validationError {
		t.Fatalf("invalid temporal = %#v, %t", cliErr, handled)
	}
}

func TestExecuteHistory(t *testing.T) {
	original := userConfigDir
	t.Cleanup(func() { userConfigDir = original })
	parsed, err := parseCommand([]string{"attach-evidence", "--assertion-id", "missing", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"})
	if err.code != "" {
		t.Fatalf("parse: %v", err)
	}
	userConfigDir = func() (string, error) { return t.TempDir(), nil }
	if _, cliErr := executeHistory(context.Background(), parsed); cliErr.code != notFoundError {
		t.Fatalf("not found = %#v", cliErr)
	}
	userConfigDir = func() (string, error) { return "", errors.New("config") }
	if _, cliErr := executeHistory(context.Background(), parsed); cliErr.code != storageError {
		t.Fatalf("store failure = %#v", cliErr)
	}
}
