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

// CreateService は作成操作を実行する。
type CreateService struct {
	store domain.CreateStore
}

// HistoryService は履歴mutationを実行する。
type HistoryService struct{ store domain.HistoryStore }

// NewHistoryService は履歴mutationサービスを作る。
func NewHistoryService(store domain.HistoryStore) HistoryService { return HistoryService{store: store} }

// AttachEvidence は根拠を追加する。
func (s HistoryService) AttachEvidence(ctx context.Context, request domain.AttachEvidenceRequest) (domain.AttachEvidenceResult, error) {
	return s.store.AttachEvidence(ctx, request)
}

// Revise は新しいrevisionを追加する。
func (s HistoryService) Revise(ctx context.Context, request domain.ReviseRequest) (domain.ReviseResult, error) {
	return s.store.ReviseAssertion(ctx, request)
}

// Supersede は置換Relationを追加する。
func (s HistoryService) Supersede(ctx context.Context, request domain.SupersedeRequest) (domain.SupersedeResult, error) {
	return s.store.Supersede(ctx, request)
}

// NewCreateService は作成操作のサービスを作る。
func NewCreateService(store domain.CreateStore) CreateService {
	return CreateService{store: store}
}

// Create はAssertion初版を作成する。
func (s CreateService) Create(ctx context.Context, request domain.CreateRequest) (domain.CreateResult, error) {
	return s.store.CreateAssertion(ctx, request)
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

// SearchRelated は検索起点に接続するRelationを取得する。
func (s RetrievalService) SearchRelated(ctx context.Context, kind string, id string, relationTypes []string) ([]domain.RelatedResult, error) {
	return s.store.SearchRelated(ctx, kind, id, relationTypes)
}

// SearchContradictions は矛盾候補を取得する。
func (s RetrievalService) SearchContradictions(ctx context.Context, assertionID *string, concept *string) ([]domain.ContradictionResult, error) {
	return s.store.SearchContradictions(ctx, assertionID, concept)
}

// SearchTemporal は時点情報を持つ現行Assertionを取得する。
func (s RetrievalService) SearchTemporal(ctx context.Context, concept *string, scope []domain.Scope, filter domain.TemporalSearchFilter) ([]domain.TemporalSearchResult, error) {
	return s.store.SearchTemporal(ctx, concept, scope, filter)
}
