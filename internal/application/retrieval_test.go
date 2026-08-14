package application

import (
	"context"
	"testing"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

type retrievalStoreStub struct {
	receivedContext *context.Context
}

func (store retrievalStoreStub) SearchText(ctx context.Context, _ string) ([]domain.AssertionSummary, error) {
	store.receive(ctx)

	return []domain.AssertionSummary{{ID: "asrt_01"}}, nil
}

func (store retrievalStoreStub) SearchConcept(ctx context.Context, _ string) (domain.ConceptSearchResult, error) {
	store.receive(ctx)

	return domain.ConceptSearchResult{Concept: &domain.Concept{ID: "cpt_01"}}, nil
}

func (store retrievalStoreStub) GetAssertion(ctx context.Context, _ string) (domain.AssertionDetail, error) {
	store.receive(ctx)

	return domain.AssertionDetail{ID: "asrt_01"}, nil
}

func (store retrievalStoreStub) GetEvidence(ctx context.Context, _ string) (domain.EvidenceResult, error) {
	store.receive(ctx)

	return domain.EvidenceResult{AssertionID: "asrt_01"}, nil
}

func (store retrievalStoreStub) receive(ctx context.Context) {
	if store.receivedContext != nil {
		*store.receivedContext = ctx
	}
}

func TestRetrievalService(t *testing.T) {
	tests := []struct {
		name string
		call func(*testing.T, RetrievalService, context.Context)
	}{
		{
			name: "SearchText",
			call: func(t *testing.T, service RetrievalService, ctx context.Context) {
				results, err := service.SearchText(ctx, "channel")
				if err != nil || results[0].ID != "asrt_01" {
					t.Fatalf("SearchText() = %#v, %v", results, err)
				}
			},
		},
		{
			name: "SearchConcept",
			call: func(t *testing.T, service RetrievalService, ctx context.Context) {
				result, err := service.SearchConcept(ctx, "channel")
				if err != nil || result.Concept.ID != "cpt_01" {
					t.Fatalf("SearchConcept() = %#v, %v", result, err)
				}
			},
		},
		{
			name: "Get",
			call: func(t *testing.T, service RetrievalService, ctx context.Context) {
				result, err := service.Get(ctx, "asrt_01")
				if err != nil || result.ID != "asrt_01" {
					t.Fatalf("Get() = %#v, %v", result, err)
				}
			},
		},
		{
			name: "GetEvidence",
			call: func(t *testing.T, service RetrievalService, ctx context.Context) {
				result, err := service.GetEvidence(ctx, "asrt_01")
				if err != nil || result.AssertionID != "asrt_01" {
					t.Fatalf("GetEvidence() = %#v, %v", result, err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var received context.Context
			service := NewRetrievalService(retrievalStoreStub{receivedContext: &received})
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			tt.call(t, service, ctx)
			if received != ctx {
				t.Fatalf("%sのContextがStoreへ伝播しません", tt.name)
			}
		})
	}
}
