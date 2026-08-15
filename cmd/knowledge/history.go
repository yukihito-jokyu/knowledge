package main

import (
	"context"
	"errors"

	"github.com/yukihito-jokyu/knowledge/internal/application"
	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

// executeHistory は既定Storeを開いて履歴mutationを実行する。
func executeHistory(ctx context.Context, parsed command) (any, cliError) {
	store, err := openDefaultRetrievalStore(ctx)
	if err != nil {
		return nil, cliError{
			code:    storageError,
			message: "Knowledge Storeを開けません",
		}
	}
	defer store.Close()
	data, executionError, _ := executeHistoryWithStore(ctx, parsed, store)

	return data, executionError
}

func executeHistoryWithStore(ctx context.Context, parsed command, store domain.HistoryStore) (any, cliError, bool) {
	values := optionValues(parsed.options)
	service := application.NewHistoryService(store)
	switch parsed.operation {
	case "attach-evidence":
		result, err := service.AttachEvidence(ctx, domain.AttachEvidenceRequest{
			AssertionID: values["assertion-id"][0],
			Evidence: domain.CreateEvidence{
				Kind:       values["evidence-kind"][0],
				RawText:    values["evidence-text"][0],
				ObservedAt: canonicalCreateTime(values["evidence-observed-at"][0]),
			},
		})
		if err != nil {
			return nil, historyCLIError(err), true
		}

		return map[string]any{
			"assertion_id": result.AssertionID,
			"evidence_id":  result.EvidenceID,
		}, cliError{}, true
	case "revise":
		temporal, cliErr := createTemporal(values)
		if cliErr.code != "" {
			return nil, cliErr, true
		}
		result, err := service.Revise(ctx, domain.ReviseRequest{
			AssertionID:    values["assertion-id"][0],
			NormalizedText: values["normalized-text"][0],
			Scope:          scopeValues(values),
			Temporal:       temporal,
		})
		if err != nil {
			return nil, historyCLIError(err), true
		}

		return map[string]any{
			"assertion_id":      result.AssertionID,
			"previous_revision": result.PreviousRevision,
			"revision":          result.Revision,
		}, cliError{}, true
	case "supersede":
		result, err := service.Supersede(ctx, domain.SupersedeRequest{
			SupersededAssertionID:  values["superseded-assertion-id"][0],
			ReplacementAssertionID: values["replacement-assertion-id"][0],
		})
		if err != nil {
			return nil, historyCLIError(err), true
		}

		return map[string]any{
			"relation_id":              result.RelationID,
			"relation_type":            "supersedes",
			"superseded_assertion_id":  result.SupersededAssertionID,
			"replacement_assertion_id": result.ReplacementAssertionID,
		}, cliError{}, true
	default:
		return nil, cliError{}, false
	}
}

func historyCLIError(err error) cliError {
	if errors.Is(err, domain.ErrAssertionNotFound) {
		return cliError{
			code:    notFoundError,
			message: "Assertionが見つかりません",
		}
	}
	if errors.Is(err, domain.ErrMutationConflict) {
		return cliError{
			code:    conflictError,
			message: "履歴更新が既存データと衝突しました",
		}
	}

	return cliError{
		code:    storageError,
		message: "Knowledge Storeの履歴更新に失敗しました",
	}
}
