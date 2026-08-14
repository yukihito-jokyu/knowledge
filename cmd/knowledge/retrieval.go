package main

import (
	"context"
	"errors"

	"github.com/yukihito-jokyu/knowledge/internal/application"
	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

func executeRetrievalWithStore(ctx context.Context, parsed command, store domain.RetrievalStore) (any, cliError, bool) {
	service := application.NewRetrievalService(store)
	values := optionValues(parsed.options)
	var data any
	var err error
	switch parsed.operation {
	case "search-text":
		data, err = service.SearchText(ctx, values["query"][0])
	case "search-concept":
		data, err = service.SearchConcept(ctx, values["concept"][0])
	case "get":
		data, err = service.Get(ctx, values["assertion-id"][0])
	case "get-evidence":
		data, err = service.GetEvidence(ctx, values["assertion-id"][0])
	default:
		return nil, cliError{}, false
	}
	if err != nil {
		return nil, retrievalCLIError(err), true
	}

	return retrievalResponse(data), cliError{}, true
}

func retrievalCLIError(err error) cliError {
	if errors.Is(err, domain.ErrAssertionNotFound) {
		return cliError{
			code:    notFoundError,
			message: "Assertionが見つかりません",
		}
	}

	return cliError{
		code:    storageError,
		message: "Knowledge Storeの読取に失敗しました",
	}
}

func retrievalResponse(data any) any {
	switch value := data.(type) {
	case []domain.AssertionSummary:
		return map[string]any{
			"results": summariesResponse(value),
		}
	case domain.ConceptSearchResult:
		return map[string]any{
			"concept": conceptResponse(value.Concept),
			"results": conceptSummariesResponse(value.Results),
		}
	case domain.AssertionDetail:
		return assertionResponse(value)
	case domain.EvidenceResult:
		evidence := make([]map[string]any, 0, len(value.Evidence))
		for _, entry := range value.Evidence {
			evidence = append(evidence, map[string]any{
				"evidence_id": entry.ID,
				"kind":        entry.Kind,
				"raw_text":    entry.RawText,
				"observed_at": entry.ObservedAt,
			})
		}

		return map[string]any{
			"assertion_id": value.AssertionID,
			"evidence":     evidence,
		}
	default:
		return nil
	}
}

func summariesResponse(values []domain.AssertionSummary) []map[string]any {
	results := make([]map[string]any, 0, len(values))
	for _, value := range values {
		results = append(results, map[string]any{
			"assertion_id":    value.ID,
			"normalized_text": value.NormalizedText,
			"revision":        value.Revision,
			"concepts":        conceptsResponse(value.Concepts),
			"scope":           scopeResponse(value.Scope),
			"matched_fields":  value.MatchedFields,
		})
	}

	return results
}

func conceptSummariesResponse(values []domain.AssertionSummary) []map[string]any {
	results := make([]map[string]any, 0, len(values))
	for _, value := range values {
		results = append(results, map[string]any{
			"assertion_id":    value.ID,
			"normalized_text": value.NormalizedText,
			"revision":        value.Revision,
			"scope":           scopeResponse(value.Scope),
		})
	}

	return results
}

func conceptsResponse(values []domain.Concept) []map[string]any {
	results := make([]map[string]any, 0, len(values))
	for _, value := range values {
		results = append(results, map[string]any{
			"concept_id": value.ID,
			"name":       value.Name,
		})
	}

	return results
}

func scopeResponse(values []domain.Scope) []map[string]any {
	results := make([]map[string]any, 0, len(values))
	for _, value := range values {
		results = append(results, map[string]any{
			"key":   value.Key,
			"value": value.Value,
		})
	}

	return results
}

func conceptResponse(value *domain.Concept) any {
	if value == nil {
		return nil
	}

	return map[string]any{
		"concept_id": value.ID,
		"name":       value.Name,
	}
}

func assertionResponse(value domain.AssertionDetail) map[string]any {
	return map[string]any{
		"assertion_id":     value.ID,
		"current_revision": value.CurrentRevision,
		"revisions":        revisionsResponse(value.Revisions),
		"concepts":         conceptDetailsResponse(value.Concepts),
		"aliases":          assertionAliasesResponse(value.Aliases),
	}
}

func revisionsResponse(values []domain.Revision) []map[string]any {
	results := make([]map[string]any, 0, len(values))
	for _, value := range values {
		results = append(results, map[string]any{
			"revision":        value.Number,
			"normalized_text": value.NormalizedText,
			"scope":           scopeResponse(value.Scope),
			"temporal":        temporalResponse(value.Temporal),
		})
	}

	return results
}

func temporalResponse(value *domain.Temporal) any {
	if value == nil {
		return nil
	}

	return map[string]any{
		"valid_from":    value.ValidFrom,
		"valid_until":   value.ValidUntil,
		"version_scope": value.VersionScope,
		"observed_at":   value.ObservedAt,
		"last_verified": value.LastVerified,
	}
}

func conceptDetailsResponse(values []domain.ConceptDetail) []map[string]any {
	results := make([]map[string]any, 0, len(values))
	for _, value := range values {
		results = append(results, map[string]any{
			"concept_id": value.ID,
			"name":       value.Name,
			"aliases":    value.Aliases,
		})
	}

	return results
}

func assertionAliasesResponse(values []domain.AssertionAlias) []map[string]any {
	results := make([]map[string]any, 0, len(values))
	for _, value := range values {
		results = append(results, map[string]any{
			"kind":  value.Kind,
			"value": value.Value,
		})
	}

	return results
}
