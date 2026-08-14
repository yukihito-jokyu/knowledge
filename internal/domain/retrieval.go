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

// Concept は検索アンカーとなるConceptを表す。
type Concept struct {
	ID   string
	Name string
}

// Scope はAssertion revisionの適用条件を表す。
type Scope struct {
	Key   string
	Value string
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
}
