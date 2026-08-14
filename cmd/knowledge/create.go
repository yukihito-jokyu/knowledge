package main

import (
	"context"
	"errors"
	"time"

	"github.com/yukihito-jokyu/knowledge/internal/application"
	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

// executeCreate は既定Storeを開いてcreate操作を実行する。
func executeCreate(ctx context.Context, parsed command) (any, cliError) {
	store, err := openDefaultRetrievalStore(ctx)
	if err != nil {
		return nil, cliError{
			code:    storageError,
			message: "Knowledge Storeを開けません",
		}
	}
	defer store.Close()

	data, executionError, _ := executeCreateWithStore(ctx, parsed, store)

	return data, executionError
}

func executeCreateWithStore(ctx context.Context, parsed command, store domain.CreateStore) (any, cliError, bool) {
	if parsed.operation != "create" {
		return nil, cliError{}, false
	}
	request, err := createRequest(parsed.options)
	if err.code != "" {
		return nil, err, true
	}
	result, createErr := application.NewCreateService(store).Create(ctx, request)
	if createErr != nil {
		return nil, createCLIError(createErr), true
	}

	return map[string]any{
		"assertion_id": result.AssertionID,
		"revision":     result.Revision,
		"evidence_ids": result.EvidenceIDs,
		"concepts":     conceptsResponse(result.Concepts),
		"relation_ids": result.RelationIDs,
	}, cliError{}, true
}

func createRequest(options []option) (domain.CreateRequest, cliError) {
	values := optionValues(options)
	temporal, temporalErr := createTemporal(values)
	if temporalErr.code != "" {
		return domain.CreateRequest{}, temporalErr
	}
	request := domain.CreateRequest{
		NormalizedText: values["normalized-text"][0],
		Scope:          scopeValues(values),
		Aliases:        pairedAliases(values),
		Relations:      pairedRelations(values),
		Evidence:       pairedEvidence(values),
		Temporal:       temporal,
	}
	for _, entry := range options {
		if entry.name == "concept" {
			request.Concepts = append(request.Concepts, domain.CreateConcept{
				Name:    entry.value,
				Aliases: make([]string, 0),
			})

			continue
		}
		if entry.name == "concept-alias" {
			request.Concepts[len(request.Concepts)-1].Aliases = append(request.Concepts[len(request.Concepts)-1].Aliases, entry.value)
		}
	}

	return request, cliError{}
}

func pairedAliases(values map[string][]string) []domain.AssertionAlias {
	aliases := make([]domain.AssertionAlias, 0, len(values["alias-kind"]))
	for index, kind := range values["alias-kind"] {
		aliases = append(aliases, domain.AssertionAlias{
			Kind:  kind,
			Value: values["alias-value"][index],
		})
	}

	return aliases
}

func pairedRelations(values map[string][]string) []domain.CreateRelation {
	relations := make([]domain.CreateRelation, 0, len(values["relation-type"]))
	for index, kind := range values["relation-type"] {
		relations = append(relations, domain.CreateRelation{
			Type:       kind,
			TargetKind: values["relation-target-kind"][index],
			TargetID:   values["relation-target-id"][index],
		})
	}

	return relations
}

func pairedEvidence(values map[string][]string) []domain.CreateEvidence {
	evidence := make([]domain.CreateEvidence, 0, len(values["evidence-kind"]))
	for index, kind := range values["evidence-kind"] {
		evidence = append(evidence, domain.CreateEvidence{
			Kind:       kind,
			RawText:    values["evidence-text"][index],
			ObservedAt: canonicalCreateTime(values["evidence-observed-at"][index]),
		})
	}

	return evidence
}

func createTemporal(values map[string][]string) (*domain.Temporal, cliError) {
	if len(values["valid-from"])+len(values["valid-until"])+len(values["version-scope"])+len(values["observed-at"])+len(values["last-verified"]) == 0 {
		return nil, cliError{}
	}
	temporal := &domain.Temporal{}
	for _, entry := range []struct {
		name   string
		target **string
	}{
		{
			name:   "valid-from",
			target: &temporal.ValidFrom,
		},
		{
			name:   "valid-until",
			target: &temporal.ValidUntil,
		},
		{
			name:   "observed-at",
			target: &temporal.ObservedAt,
		},
		{
			name:   "last-verified",
			target: &temporal.LastVerified,
		},
	} {
		value := optionalValue(values, entry.name)
		if value != nil {
			parsed, parseErr := time.Parse(time.RFC3339, *value)
			if parseErr != nil {
				return nil, validationFailure(entry.name, "RFC 3339 UTC時刻を指定してください")
			}
			canonical := parsed.UTC().Format(fixedTemporalTimestampLayout)
			*entry.target = &canonical
		}
	}
	temporal.VersionScope = optionalValue(values, "version-scope")

	return temporal, cliError{}
}

func canonicalCreateTime(value string) string {
	parsed, _ := time.Parse(time.RFC3339, value)

	return parsed.UTC().Format(fixedTemporalTimestampLayout)
}

func createCLIError(err error) cliError {
	if errors.Is(err, domain.ErrCreateConflict) {
		return cliError{
			code:    conflictError,
			message: "作成対象が既存データと衝突しました",
		}
	}
	if errors.Is(err, domain.ErrCreateRelationTargetNotFound) {
		return cliError{
			code:    notFoundError,
			message: "Relationの参照先が見つかりません",
		}
	}

	return cliError{
		code:    storageError,
		message: "Knowledge Storeへの作成に失敗しました",
	}
}
