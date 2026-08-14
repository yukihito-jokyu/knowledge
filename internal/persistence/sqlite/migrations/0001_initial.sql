-- このStoreへ正常に適用済みのmigration versionを記録する。
CREATE TABLE schema_migrations (
  version INTEGER PRIMARY KEY,
  applied_at TEXT NOT NULL
);

-- 安定したAssertion識別子と現行revision番号を保持する。
CREATE TABLE assertions (
  assertion_id TEXT PRIMARY KEY,
  current_revision INTEGER NOT NULL CHECK (current_revision >= 1),
  created_at TEXT NOT NULL
);

-- Assertionの正規化済み本文を不変のrevisionとして保持する。
CREATE TABLE assertion_revisions (
  assertion_id TEXT NOT NULL,
  revision INTEGER NOT NULL CHECK (revision >= 1),
  normalized_text TEXT NOT NULL CHECK (length(trim(normalized_text)) > 0),
  created_at TEXT NOT NULL,
  PRIMARY KEY (assertion_id, revision),
  FOREIGN KEY (assertion_id) REFERENCES assertions(assertion_id)
);

-- 特定のAssertion revisionに属するScopeのkey/valueを保持する。
CREATE TABLE revision_scopes (
  assertion_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  scope_key TEXT NOT NULL CHECK (length(trim(scope_key)) > 0),
  scope_value TEXT NOT NULL CHECK (length(trim(scope_value)) > 0),
  PRIMARY KEY (assertion_id, revision, scope_key),
  FOREIGN KEY (assertion_id, revision) REFERENCES assertion_revisions(assertion_id, revision)
);

-- Assertionに紐付く不変のEvidenceを保持する。
CREATE TABLE evidence (
  evidence_id TEXT PRIMARY KEY,
  assertion_id TEXT NOT NULL,
  kind TEXT NOT NULL CHECK (kind IN ('user_explanation', 'user_reasoning', 'user_code', 'self_report', 'concept_recognition', 'correction')),
  raw_text TEXT NOT NULL CHECK (length(trim(raw_text)) > 0),
  observed_at TEXT NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY (assertion_id) REFERENCES assertions(assertion_id)
);

-- 検索アンカーに使うConceptの正規名を保持する。
CREATE TABLE concepts (
  concept_id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE CHECK (length(trim(name)) > 0),
  created_at TEXT NOT NULL
);

-- Conceptの正規名とAliasを横断して一意にする検索語台帳を保持する。
CREATE TABLE concept_terms (
  term TEXT PRIMARY KEY CHECK (length(trim(term)) > 0),
  concept_id TEXT NOT NULL,
  term_kind TEXT NOT NULL CHECK (term_kind IN ('name', 'alias')),
  FOREIGN KEY (concept_id) REFERENCES concepts(concept_id)
);

-- Conceptに属するAliasを保持する。
CREATE TABLE concept_aliases (
  alias TEXT PRIMARY KEY CHECK (length(trim(alias)) > 0),
  concept_id TEXT NOT NULL,
  FOREIGN KEY (concept_id) REFERENCES concepts(concept_id)
);

-- 現行AssertionとConceptの多対多対応を保持する。
CREATE TABLE assertion_concepts (
  assertion_id TEXT NOT NULL,
  concept_id TEXT NOT NULL,
  PRIMARY KEY (assertion_id, concept_id),
  FOREIGN KEY (assertion_id) REFERENCES assertions(assertion_id),
  FOREIGN KEY (concept_id) REFERENCES concepts(concept_id)
);

-- Assertionに付随するAPI名・Identifierの検索語を保持する。
CREATE TABLE assertion_aliases (
  assertion_id TEXT NOT NULL,
  alias_kind TEXT NOT NULL CHECK (alias_kind IN ('api_name', 'identifier')),
  alias_value TEXT NOT NULL CHECK (length(trim(alias_value)) > 0),
  PRIMARY KEY (assertion_id, alias_kind, alias_value),
  FOREIGN KEY (assertion_id) REFERENCES assertions(assertion_id)
);

-- Assertion revisionに任意で紐付く時点・有効期間情報を保持する。
CREATE TABLE temporal_metadata (
  assertion_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  valid_from TEXT,
  valid_until TEXT,
  version_scope TEXT,
  observed_at TEXT,
  last_verified TEXT,
  PRIMARY KEY (assertion_id, revision),
  CHECK (valid_from IS NULL OR valid_until IS NULL OR valid_from <= valid_until),
  FOREIGN KEY (assertion_id, revision) REFERENCES assertion_revisions(assertion_id, revision)
);

-- AssertionまたはConcept endpoint間の明示済みRelationを保持する。
CREATE TABLE relations (
  relation_id TEXT PRIMARY KEY,
  source_kind TEXT NOT NULL CHECK (source_kind IN ('assertion', 'concept')),
  source_id TEXT NOT NULL,
  relation_type TEXT NOT NULL CHECK (relation_type IN ('related_to', 'prerequisite', 'causes', 'contributes_to', 'contradicts', 'supersedes')),
  target_kind TEXT NOT NULL CHECK (target_kind IN ('assertion', 'concept')),
  target_id TEXT NOT NULL,
  created_at TEXT NOT NULL,
  CHECK (NOT (source_kind = target_kind AND source_id = target_id)),
  CHECK (relation_type <> 'supersedes' OR (source_kind = 'assertion' AND target_kind = 'assertion')),
  CHECK (relation_type <> 'contradicts' OR (source_kind = 'assertion' AND target_kind = 'assertion')),
  UNIQUE (source_kind, source_id, relation_type, target_kind, target_id)
);

-- 現行Assertion内容の派生FTS5字句Indexを提供する。
CREATE VIRTUAL TABLE assertion_lexical_index USING fts5(
  assertion_id UNINDEXED,
  normalized_text,
  concept_name,
  concept_alias,
  scope_key,
  scope_value,
  assertion_alias
);

CREATE INDEX idx_scopes_key_value ON revision_scopes(scope_key, scope_value);
CREATE INDEX idx_evidence_assertion ON evidence(assertion_id, observed_at);
CREATE INDEX idx_aliases_concept ON concept_aliases(concept_id);
CREATE INDEX idx_concept_terms_concept ON concept_terms(concept_id);
CREATE INDEX idx_assertion_concepts_concept ON assertion_concepts(concept_id, assertion_id);
CREATE INDEX idx_assertion_aliases_value ON assertion_aliases(alias_value, assertion_id);
CREATE INDEX idx_relations_source ON relations(source_kind, source_id, relation_type);
CREATE INDEX idx_relations_target ON relations(target_kind, target_id, relation_type);
CREATE INDEX idx_temporal_window ON temporal_metadata(valid_from, valid_until);
