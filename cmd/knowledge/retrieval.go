package main

import (
	"context"
	"errors"
	"time"

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
	case "search-related":
		data, err = service.SearchRelated(ctx, values["seed-kind"][0], values["seed-id"][0], values["relation-type"])
	case "search-contradictions":
		data, err = service.SearchContradictions(ctx, optionalValue(values, "assertion-id"), optionalValue(values, "concept"))
	case "search-temporal":
		filter, filterError := temporalSearchFilter(values)
		if filterError.code != "" {
			return nil, filterError, true
		}
		data, err = service.SearchTemporal(ctx, optionalValue(values, "concept"), scopeValues(values), filter)
	default:
		return nil, cliError{}, false
	}
	if err != nil {
		return nil, retrievalCLIError(err), true
	}

	return retrievalResponse(data), cliError{}, true
}

const fixedTemporalTimestampLayout = "2006-01-02T15:04:05.000000000Z"

func temporalSearchFilter(values map[string][]string) (domain.TemporalSearchFilter, cliError) {
	filter := domain.TemporalSearchFilter{}
	for _, entry := range []struct {
		name   string
		target **string
	}{
		{
			name:   "at",
			target: &filter.At,
		},
		{
			name:   "valid-from",
			target: &filter.ValidFrom,
		},
		{
			name:   "valid-until",
			target: &filter.ValidUntil,
		},
	} {
		value := optionalValue(values, entry.name)
		if value == nil {
			continue
		}
		parsed, err := time.Parse(time.RFC3339, *value)
		if err != nil {
			return domain.TemporalSearchFilter{}, validationFailure(entry.name, "RFC 3339 UTC時刻を指定してください")
		}
		_, offset := parsed.Zone()
		if offset != 0 {
			return domain.TemporalSearchFilter{}, validationFailure(entry.name, "UTC時刻を指定してください")
		}
		canonical := parsed.UTC().Format(fixedTemporalTimestampLayout)
		*entry.target = &canonical
	}

	return filter, cliError{}
}

func retrievalCLIError(err error) cliError {
	if errors.Is(err, domain.ErrRelationSeedNotFound) {
		return cliError{
			code:    notFoundError,
			message: "検索起点が見つかりません",
		}
	}
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
			item := map[string]any{
				"evidence_id": entry.ID,
				"kind":        entry.Kind,
				"raw_text":    entry.RawText,
				"observed_at": entry.ObservedAt,
			}
			if entry.Temporal != nil {
				item["temporal"] = temporalResponse(entry.Temporal)
			}
			evidence = append(evidence, item)
		}

		return map[string]any{
			"assertion_id": value.AssertionID,
			"evidence":     evidence,
		}
	case []domain.RelatedResult:
		results := make([]map[string]any, 0, len(value))
		for _, entry := range value {
			results = append(results, relatedResponse(entry))
		}

		return map[string]any{"results": results}
	case []domain.ContradictionResult:
		results := make([]map[string]any, 0, len(value))
		for _, entry := range value {
			results = append(results, contradictionResponse(entry))
		}

		return map[string]any{"results": results}
	case []domain.TemporalSearchResult:
		results := make([]map[string]any, 0, len(value))
		for _, entry := range value {
			results = append(results, map[string]any{
				"assertion_id":    entry.AssertionID,
				"normalized_text": entry.NormalizedText,
				"temporal":        temporalResponse(&entry.Temporal),
			})
		}

		return map[string]any{"results": results}
	default:
		return nil
	}
}

func optionalValue(values map[string][]string, name string) *string {
	if len(values[name]) == 0 {
		return nil
	}

	return &values[name][0]
}

func scopeValues(values map[string][]string) []domain.Scope {
	keys := values["scope-key"]
	entries := values["scope-value"]
	scope := make([]domain.Scope, 0, len(keys))
	for index, key := range keys {
		scope = append(scope, domain.Scope{
			Key:   key,
			Value: entries[index],
		})
	}

	return scope
}

func relatedResponse(value domain.RelatedResult) map[string]any {
	return map[string]any{
		"relation_id":   value.RelationID,
		"relation_type": value.RelationType,
		"direction":     value.Direction,
		"target":        relationTargetResponse(value.Target),
	}
}

func contradictionResponse(value domain.ContradictionResult) map[string]any {
	return map[string]any{
		"relation_id": value.RelationID,
		"direction":   value.Direction,
		"seed": map[string]any{
			"kind": "assertion",
			"id":   value.SeedID,
		},
		"target": relationTargetResponse(value.Target),
	}
}

func relationTargetResponse(value domain.RelationTarget) map[string]any {
	result := map[string]any{
		"kind": value.Kind,
		"id":   value.ID,
	}
	if value.NormalizedText != nil {
		result["normalized_text"] = *value.NormalizedText
	}

	return result
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
