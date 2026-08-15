-- Evidenceごとの任意のTemporal metadataを保持する。
CREATE TABLE evidence_temporal_metadata (
  evidence_id TEXT PRIMARY KEY,
  valid_from TEXT,
  valid_until TEXT,
  version_scope TEXT,
  observed_at TEXT,
  last_verified TEXT,
  CHECK (valid_from IS NULL OR valid_until IS NULL OR valid_from <= valid_until),
  FOREIGN KEY (evidence_id) REFERENCES evidence(evidence_id)
);

CREATE INDEX idx_evidence_temporal_window ON evidence_temporal_metadata(valid_from, valid_until);
