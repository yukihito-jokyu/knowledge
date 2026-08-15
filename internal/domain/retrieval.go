// Package domain はKnowledge CLIの意味的に中立な値と永続化portを定義する。
package domain

import (
	"context"
	"errors"
)

// ErrAssertionNotFound は指定されたAssertionが存在しないことを示す。
var ErrAssertionNotFound = errors.New("assertion not found")

// ErrRelationSeedNotFound はRelation検索起点が存在しないことを示す。
var ErrRelationSeedNotFound = errors.New("relation seed not found")

// ErrCreateConflict は作成対象が既存データと衝突したことを示す。
var ErrCreateConflict = errors.New("create conflict")

// ErrCreateRelationTargetNotFound はRelationの参照先が存在しないことを示す。
var ErrCreateRelationTargetNotFound = errors.New("create relation target not found")

// ErrMutationConflict は履歴mutationが既存データと衝突したことを示す。
var ErrMutationConflict = errors.New("mutation conflict")

// Concept は検索アンカーとなるConceptを表す。
type Concept struct {
	ID   string
	Name string
}

// Scope はAssertion revisionの適用条件を表す。
type Scope struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AssertionSummary は現行Assertionの検索結果を表す。
type AssertionSummary struct {
	ID             string
	NormalizedText string
	Revision       int
	Concepts       []Concept
	Scope          []Scope
	MatchedFields  []string
}

// ConceptSearchResult はConcept検索の結果を表す。
type ConceptSearchResult struct {
	Concept *Concept
	Results []AssertionSummary
}

// Temporal はrevisionに紐付く時点情報を表す。
type Temporal struct {
	ValidFrom    *string
	ValidUntil   *string
	VersionScope *string
	ObservedAt   *string
	LastVerified *string
}

// TemporalSearchResult は時点情報を持つ現行Assertionの検索結果を表す。
type TemporalSearchResult struct {
	AssertionID    string
	NormalizedText string
	Temporal       Temporal
}

// TemporalSearchFilter は時点または有効期間による検索条件を表す。
type TemporalSearchFilter struct {
	At         *string
	ValidFrom  *string
	ValidUntil *string
}

// Revision はAssertionの不変なrevisionを表す。
type Revision struct {
	Number         int
	NormalizedText string
	Scope          []Scope
	Temporal       *Temporal
}

// ConceptDetail はConceptとAliasを表す。
type ConceptDetail struct {
	ID      string
	Name    string
	Aliases []string
}

// AssertionAlias はAssertionに直接付与した検索語を表す。
type AssertionAlias struct {
	Kind  string
	Value string
}

// AssertionDetail はAssertionの取得結果を表す。
type AssertionDetail struct {
	ID              string
	CurrentRevision int
	Revisions       []Revision
	Concepts        []ConceptDetail
	Aliases         []AssertionAlias
}

// Evidence はAssertionに紐付く根拠記録を表す。
type Evidence struct {
	ID         string
	Kind       string
	RawText    string
	ObservedAt string
}

// EvidenceResult はEvidence取得結果を表す。
type EvidenceResult struct {
	AssertionID string
	Evidence    []Evidence
}

// RelationTarget はRelation検索で検索起点の反対側を表す。
type RelationTarget struct {
	Kind           string
	ID             string
	NormalizedText *string
}

// RelatedResult は検索起点に接続するRelationを表す。
type RelatedResult struct {
	RelationID   string
	RelationType string
	Direction    string
	Target       RelationTarget
}

// ContradictionResult は矛盾候補の検索起点と相手を表す。
type ContradictionResult struct {
	RelationID string
	Direction  string
	SeedID     string
	Target     RelationTarget
}

// RetrievalStore は読取操作が必要とする永続化portである。
type RetrievalStore interface {
	SearchText(context.Context, string) ([]AssertionSummary, error)
	SearchConcept(context.Context, string) (ConceptSearchResult, error)
	GetAssertion(context.Context, string) (AssertionDetail, error)
	GetEvidence(context.Context, string) (EvidenceResult, error)
	SearchRelated(context.Context, string, string, []string) ([]RelatedResult, error)
	SearchContradictions(context.Context, *string, *string) ([]ContradictionResult, error)
	SearchTemporal(context.Context, *string, []Scope, TemporalSearchFilter) ([]TemporalSearchResult, error)
}

// CreateRequest はAssertion初版と付随データの作成入力を表す。
type CreateRequest struct {
	NormalizedText string
	Scope          []Scope
	Concepts       []CreateConcept
	Aliases        []AssertionAlias
	Relations      []CreateRelation
	Evidence       []CreateEvidence
	Temporal       *Temporal
}

type CreateConcept struct {
	Name    string
	Aliases []string
}

type CreateRelation struct {
	Type       string
	TargetKind string
	TargetID   string
}

type CreateEvidence struct {
	Kind       string
	RawText    string
	ObservedAt string
}

// CreateResult はcommit済みの作成結果を表す。
type CreateResult struct {
	AssertionID string
	Revision    int
	EvidenceIDs []string
	Concepts    []Concept
	RelationIDs []string
}

// CreateStore は作成操作が必要とする永続化portである。
type CreateStore interface {
	CreateAssertion(context.Context, CreateRequest) (CreateResult, error)
}

// AttachEvidenceRequest は既存Assertionへ追加する根拠を表す。
type AttachEvidenceRequest struct {
	AssertionID string
	Evidence    CreateEvidence
}

// AttachEvidenceResult はcommit済みの根拠追加結果を表す。
type AttachEvidenceResult struct {
	AssertionID string
	EvidenceID  string
}

// ReviseRequest はAssertionの新しいrevisionを表す。
type ReviseRequest struct {
	AssertionID    string
	NormalizedText string
	Scope          []Scope
	Temporal       *Temporal
}

// ReviseResult はcommit済みrevision更新結果を表す。
type ReviseResult struct {
	AssertionID      string
	PreviousRevision int
	Revision         int
}

// SupersedeRequest は置換Relationの追加を表す。
type SupersedeRequest struct {
	SupersededAssertionID  string
	ReplacementAssertionID string
}

// SupersedeResult はcommit済み置換Relationを表す。
type SupersedeResult struct {
	RelationID             string
	SupersededAssertionID  string
	ReplacementAssertionID string
}

// HistoryStore は履歴mutationが必要とする永続化portである。
type HistoryStore interface {
	AttachEvidence(context.Context, AttachEvidenceRequest) (AttachEvidenceResult, error)
	ReviseAssertion(context.Context, ReviseRequest) (ReviseResult, error)
	Supersede(context.Context, SupersedeRequest) (SupersedeResult, error)
}
