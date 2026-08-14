// Package application はKnowledge CLIの操作を組み立てる。
package application

import (
	"context"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

// RetrievalService は読取操作を実行する。
type RetrievalService struct {
	store domain.RetrievalStore
}

// NewRetrievalService は読取操作のサービスを作る。
func NewRetrievalService(store domain.RetrievalStore) RetrievalService {
	return RetrievalService{
		store: store,
	}
}

// SearchText は字句検索を実行する。
func (s RetrievalService) SearchText(ctx context.Context, query string) ([]domain.AssertionSummary, error) {
	return s.store.SearchText(ctx, query)
}

// SearchConcept はConcept検索を実行する。
func (s RetrievalService) SearchConcept(ctx context.Context, concept string) (domain.ConceptSearchResult, error) {
	return s.store.SearchConcept(ctx, concept)
}

// Get はAssertion詳細を取得する。
func (s RetrievalService) Get(ctx context.Context, assertionID string) (domain.AssertionDetail, error) {
	return s.store.GetAssertion(ctx, assertionID)
}

// GetEvidence はEvidence履歴を取得する。
func (s RetrievalService) GetEvidence(ctx context.Context, assertionID string) (domain.EvidenceResult, error) {
	return s.store.GetEvidence(ctx, assertionID)
}
