package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/yukihito-jokyu/knowledge/internal/domain"
)

const coverageDriverName = "sqlite-coverage-test"

var coverageQuery func(string) (driver.Rows, error)

func init() {
	sql.Register(coverageDriverName, coverageDriver{})
}

type coverageDriver struct{}

func (coverageDriver) Open(string) (driver.Conn, error) {
	return coverageConnection{}, nil
}

type coverageConnection struct{}

func (coverageConnection) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("準備できないSQL")
}

func (coverageConnection) Close() error {
	return nil
}

func (coverageConnection) Begin() (driver.Tx, error) {
	return nil, errors.New("開始できないトランザクション")
}

func (coverageConnection) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	return coverageQuery(query)
}

type coverageRows struct {
	columns  []string
	values   [][]driver.Value
	err      error
	closeErr error
	index    int
}

type unreadableMigrationFileSystem struct {
	fstest.MapFS
}

func (unreadableMigrationFileSystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("移行SQLを読めない")
}

func (r *coverageRows) Columns() []string {
	return r.columns
}

func (r *coverageRows) Close() error {
	return r.closeErr
}

func (r *coverageRows) Next(destination []driver.Value) error {
	if r.index == len(r.values) {
		if r.err != nil {
			return r.err
		}

		return io.EOF
	}
	copy(destination, r.values[r.index])
	r.index++

	return nil
}

func TestEmbeddedMigrations(t *testing.T) {
	original := migrationFileSystem
	t.Cleanup(func() { migrationFileSystem = original })

	tests := []struct {
		name       string
		fileSystem fs.FS
		wantErr    bool
	}{
		{
			name: "連番の移行を読み込む",
			fileSystem: fstest.MapFS{
				"migrations/0001_initial.sql": &fstest.MapFile{Data: []byte("CREATE TABLE one (id INTEGER)")},
				"migrations/0002_next.sql":    &fstest.MapFile{Data: []byte("CREATE TABLE two (id INTEGER)")},
			},
		},
		{
			name: "移行がない",
			fileSystem: fstest.MapFS{
				"migrations/.keep": &fstest.MapFile{Data: []byte{}},
			},
			wantErr: true,
		},
		{
			name:       "移行ディレクトリがない",
			fileSystem: fstest.MapFS{},
			wantErr:    true,
		},
		{
			name: "ファイル名が不正",
			fileSystem: fstest.MapFS{
				"migrations/initial.sql": &fstest.MapFile{Data: []byte("SELECT 1")},
			},
			wantErr: true,
		},
		{
			name: "移行SQLを読めない",
			fileSystem: unreadableMigrationFileSystem{MapFS: fstest.MapFS{
				"migrations/0001_initial.sql": &fstest.MapFile{Data: []byte("SELECT 1")},
			}},
			wantErr: true,
		},
		{
			name: "番号が連続しない",
			fileSystem: fstest.MapFS{
				"migrations/0002_initial.sql": &fstest.MapFile{Data: []byte("SELECT 1")},
			},
			wantErr: true,
		},
		{
			name: "番号が整数範囲を超える",
			fileSystem: fstest.MapFS{
				"migrations/999999999999999999999999999999999999999999_initial.sql": &fstest.MapFile{Data: []byte("SELECT 1")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrationFileSystem = tt.fileSystem
			migrations, err := embeddedMigrations()
			if (err != nil) != tt.wantErr {
				t.Fatalf("embeddedMigrations() error = %v, wantErr %t", err, tt.wantErr)
			}
			if !tt.wantErr && len(migrations) != 2 {
				t.Fatalf("migration count = %d, want 2", len(migrations))
			}
		})
	}
}

func TestSchemaStateHandlesCatalogFailures(t *testing.T) {
	tests := []struct {
		name    string
		handler func(string) (driver.Rows, error)
		wantErr bool
	}{
		{
			name: "スキーマ照会に失敗する",
			handler: func(string) (driver.Rows, error) {
				return nil, errors.New("スキーマ照会失敗")
			},
			wantErr: true,
		},
		{
			name: "スキーマ行の読み込みに失敗する",
			handler: func(string) (driver.Rows, error) {
				return &coverageRows{
					columns: []string{"name"},
					values:  [][]driver.Value{{nil}},
				}, nil
			},
			wantErr: true,
		},
		{
			name: "スキーマ走査に失敗する",
			handler: func(string) (driver.Rows, error) {
				return &coverageRows{
					columns: []string{"name"},
					err:     errors.New("スキーマ走査失敗"),
				}, nil
			},
			wantErr: true,
		},
		{
			name: "索引照会に失敗する",
			handler: func(query string) (driver.Rows, error) {
				if strings.Contains(query, "type IN") {
					return schemaObjectRows(), nil
				}

				return nil, errors.New("索引照会失敗")
			},
			wantErr: true,
		},
		{
			name: "移行履歴照会に失敗する",
			handler: func(query string) (driver.Rows, error) {
				if strings.Contains(query, "type IN") {
					return schemaObjectRows(), nil
				}
				if strings.Contains(query, "type = 'index'") {
					return &coverageRows{
						columns: []string{"name"},
						values:  [][]driver.Value{{"index"}},
					}, nil
				}

				return nil, errors.New("履歴照会失敗")
			},
			wantErr: true,
		},
		{
			name: "移行履歴の値が不正",
			handler: func(query string) (driver.Rows, error) {
				if strings.Contains(query, "type IN") {
					return schemaObjectRows(), nil
				}
				if strings.Contains(query, "type = 'index'") {
					return &coverageRows{
						columns: []string{"name"},
						values:  [][]driver.Value{{"index"}},
					}, nil
				}

				return &coverageRows{
					columns: []string{"version"},
					values:  [][]driver.Value{{"不正"}},
				}, nil
			},
			wantErr: true,
		},
		{
			name: "移行履歴の走査に失敗する",
			handler: func(query string) (driver.Rows, error) {
				if strings.Contains(query, "type IN") {
					return schemaObjectRows(), nil
				}
				if strings.Contains(query, "type = 'index'") {
					return &coverageRows{
						columns: []string{"name"},
						values:  [][]driver.Value{{"index"}},
					}, nil
				}

				return &coverageRows{
					columns: []string{"version"},
					err:     errors.New("履歴走査失敗"),
				}, nil
			},
			wantErr: true,
		},
		{
			name: "移行履歴がない",
			handler: func(query string) (driver.Rows, error) {
				if strings.Contains(query, "type IN") {
					return schemaObjectRows(), nil
				}
				if strings.Contains(query, "type = 'index'") {
					return &coverageRows{
						columns: []string{"name"},
						values:  [][]driver.Value{{"index"}},
					}, nil
				}

				return &coverageRows{columns: []string{"version"}}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coverageQuery = tt.handler
			db, err := sql.Open(coverageDriverName, "")
			if err != nil {
				t.Fatalf("sql.Open() error = %v", err)
			}
			t.Cleanup(func() { _ = db.Close() })
			_, _, err = (&Store{db: db}).schemaState(context.Background())
			if (err != nil) != tt.wantErr {
				t.Fatalf("schemaState() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

func TestOpenHandlesDatabaseOpenFailure(t *testing.T) {
	original := openDatabase
	openDatabase = func(string, string) (*sql.DB, error) { return nil, errors.New("接続失敗") }
	t.Cleanup(func() { openDatabase = original })

	if _, err := Open(context.Background(), "ignored"); err == nil {
		t.Fatal("Open() error = nil, want error")
	}
}

func TestInitializeHandlesStatesAndFailures(t *testing.T) {
	originalMigrations := loadEmbeddedMigrations
	originalState := inspectSchemaState
	t.Cleanup(func() {
		loadEmbeddedMigrations = originalMigrations
		inspectSchemaState = originalState
	})
	db, err := sql.Open(driverName, filepath.Join(t.TempDir(), "initialize.db"))
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	store := &Store{db: db}

	tests := []struct {
		name       string
		migrations []migration
		state      schemaStatus
		stateErr   error
		loadErr    error
		wantErr    bool
	}{
		{
			name:    "移行読み込み失敗",
			loadErr: errors.New("読み込み失敗"),
			wantErr: true,
		},
		{
			name:     "状態照会失敗",
			stateErr: errors.New("状態照会失敗"),
			wantErr:  true,
		},
		{
			name:    "不整合状態",
			state:   schemaInconsistent,
			wantErr: true,
		},
		{
			name:    "未知の状態",
			state:   schemaStatus(99),
			wantErr: true,
		},
		{
			name:       "空の状態",
			migrations: []migration{},
			state:      schemaEmpty,
		},
		{
			name:       "現在の状態",
			migrations: []migration{},
			state:      schemaCurrent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadEmbeddedMigrations = func() ([]migration, error) { return tt.migrations, tt.loadErr }
			inspectSchemaState = func(*Store, context.Context) (schemaStatus, []int, error) {
				return tt.state, nil, tt.stateErr
			}
			err := store.initialize(context.Background())
			if (err != nil) != tt.wantErr {
				t.Fatalf("initialize() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

func TestInitializeRejectsForeignKeyConfigurationFailure(t *testing.T) {
	db, err := sql.Open(driverName, filepath.Join(t.TempDir(), "closed.db"))
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := (&Store{db: db}).initialize(context.Background()); err == nil {
		t.Fatal("initialize() error = nil, want error")
	}
}

func TestApplyMigrationScriptsHandlesFailures(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open(driverName, filepath.Join(t.TempDir(), "apply.db"))
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	store := &Store{db: db}
	originalCommit := commitMigration
	originalGate := waitForIntegrationGate
	t.Cleanup(func() {
		commitMigration = originalCommit
		waitForIntegrationGate = originalGate
	})

	tests := []struct {
		name       string
		migrations []migration
		commitErr  error
		gateErr    error
		wantErr    bool
	}{
		{name: "移行なし"},
		{
			name: "同期失敗",
			migrations: []migration{
				{
					version: 1,
					sql:     "SELECT 1",
				},
			},
			gateErr: errors.New("同期失敗"),
			wantErr: true,
		},
		{
			name: "SQL実行失敗",
			migrations: []migration{
				{
					version: 1,
					sql:     "不正なSQL",
				},
			},
			wantErr: true,
		},
		{
			name: "履歴記録失敗",
			migrations: []migration{
				{version: 1,
					sql: "CREATE TABLE schema_migrations (version INTEGER PRIMARY KEY, applied_at TEXT)"},
				{version: 1,
					sql: "SELECT 1"},
			},
			wantErr: true,
		},
		{
			name: "確定失敗",
			migrations: []migration{
				{version: 1,
					sql: "CREATE TABLE schema_migrations (version INTEGER PRIMARY KEY, applied_at TEXT)"},
				{version: 2,
					sql: "CREATE TABLE committed_table (id INTEGER)"},
			},
			commitErr: errors.New("確定失敗"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commitMigration = func(tx *sql.Tx) error {
				if tt.commitErr != nil {
					return tt.commitErr
				}

				return tx.Commit()
			}
			waitForIntegrationGate = func(context.Context, string) error {
				return tt.gateErr
			}
			err := store.applyMigrationScripts(ctx, tt.migrations)
			if (err != nil) != tt.wantErr {
				t.Fatalf("applyMigrationScripts() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

func TestGetAssertionRejectsIntegrationGateFailure(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, filepath.Join(t.TempDir(), "gate-read.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	originalGate := waitForIntegrationGate
	t.Cleanup(func() { waitForIntegrationGate = originalGate })
	want := errors.New("read同期失敗")
	waitForIntegrationGate = func(_ context.Context, stage string) error {
		if stage == "read" {
			return want
		}

		return nil
	}
	if _, err := store.GetAssertion(ctx, "asrt_01"); !errors.Is(err, want) {
		t.Fatalf("GetAssertion() error = %v, want %v", err, want)
	}
}

func TestApplyMigrationScriptsRejectsBeginFailure(t *testing.T) {
	db, err := sql.Open(driverName, filepath.Join(t.TempDir(), "closed-transaction.db"))
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	store := &Store{db: db}
	migrations := []migration{{
		version: 1,
		sql:     "SELECT 1",
	}}
	if err := store.applyMigrationScripts(context.Background(), migrations); err == nil {
		t.Fatal("applyMigrationScripts() error = nil, want error")
	}
}

func schemaObjectRows() driver.Rows {
	values := make([][]driver.Value, 0, len(schemaObjects))
	for _, name := range schemaObjects {
		values = append(values, []driver.Value{name})
	}

	return &coverageRows{
		columns: []string{"name"},
		values:  values,
	}
}

func TestOpenAppliesV1SchemaAndRerunsWithoutChange(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "knowledge.db")

	store, err := Open(ctx, path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	assertSchemaVersion(t, store.db, 1)
	assertSchemaObjects(t, store.db)
	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	store, err = Open(ctx, path)
	if err != nil {
		t.Fatalf("second Open() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	assertSchemaVersion(t, store.db, 1)
}

func TestOpenRejectsPartialSchemaWithoutChanges(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "partial.db")
	db, err := sql.Open(driverName, path)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	if _, err := db.ExecContext(ctx, "CREATE TABLE assertions (assertion_id TEXT PRIMARY KEY)"); err != nil {
		t.Fatalf("create partial schema: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close partial database: %v", err)
	}

	_, err = Open(ctx, path)
	if !errors.Is(err, ErrInconsistentSchema) {
		t.Fatalf("Open() error = %v, want ErrInconsistentSchema", err)
	}

	db, err = sql.Open(driverName, path)
	if err != nil {
		t.Fatalf("reopen partial database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	var count int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM sqlite_schema WHERE type IN ('table', 'virtual table') AND name IN ('schema_migrations', 'assertions', 'assertion_revisions')").Scan(&count); err != nil {
		t.Fatalf("count partial schema: %v", err)
	}
	if count != 1 {
		t.Fatalf("schema object count = %d, want 1", count)
	}
}

func TestOpenRejectsInconsistentMigrationHistory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "history.db")
	store, err := Open(ctx, path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if _, err := store.db.ExecContext(ctx, "DELETE FROM schema_migrations"); err != nil {
		t.Fatalf("delete migration history: %v", err)
	}
	if _, err := store.db.ExecContext(ctx, "INSERT INTO schema_migrations (version, applied_at) VALUES (2, '2026-01-01T00:00:00Z')"); err != nil {
		t.Fatalf("insert inconsistent migration history: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	_, err = Open(ctx, path)
	if !errors.Is(err, ErrInconsistentSchema) {
		t.Fatalf("Open() error = %v, want ErrInconsistentSchema", err)
	}
}

func TestOpenRejectsMissingRequiredIndex(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "missing-index.db")
	store, err := Open(ctx, path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if _, err := store.db.ExecContext(ctx, "DROP INDEX idx_relations_target"); err != nil {
		t.Fatalf("drop required index: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	_, err = Open(ctx, path)
	if !errors.Is(err, ErrInconsistentSchema) {
		t.Fatalf("Open() error = %v, want ErrInconsistentSchema", err)
	}
}

func TestMigrationFailureRollsBackAllDDLAndVersionRecord(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db, err := sql.Open(driverName, filepath.Join(t.TempDir(), "rollback.db"))
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	store := &Store{db: db}

	err = store.applyMigrationScripts(ctx, []migration{
		{version: 1,
			sql: "CREATE TABLE rolled_back_table (id INTEGER PRIMARY KEY)"},
		{version: 2,
			sql: "this is not valid SQL"},
	})
	if err == nil {
		t.Fatal("applyMigrationScripts() error = nil, want SQL error")
	}

	var count int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM sqlite_schema WHERE name IN ('rolled_back_table', 'schema_migrations')").Scan(&count); err != nil {
		t.Fatalf("count rollback schema: %v", err)
	}
	if count != 0 {
		t.Fatalf("schema object count after rollback = %d, want 0", count)
	}
}

func TestOpenRollsBackEmbeddedMigrationWhenDDLConflicts(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "embedded-rollback.db")
	db, err := sql.Open(driverName, path)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	if _, err := db.ExecContext(ctx, "CREATE TABLE obstruction (value TEXT); CREATE INDEX idx_scopes_key_value ON obstruction(value)"); err != nil {
		t.Fatalf("create DDL obstruction: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close DDL obstruction database: %v", err)
	}

	if _, err := Open(ctx, path); err == nil {
		t.Fatal("Open() error = nil, want DDL conflict")
	}

	db, err = sql.Open(driverName, path)
	if err != nil {
		t.Fatalf("reopen database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	var count int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM sqlite_schema WHERE name IN ('schema_migrations', 'assertions', 'assertion_revisions')").Scan(&count); err != nil {
		t.Fatalf("count migration artifacts: %v", err)
	}
	if count != 0 {
		t.Fatalf("migration artifacts after DDL rollback = %d, want 0", count)
	}
}

func TestPendingMigrationsSelectsOnlyConsecutiveUnappliedVersions(t *testing.T) {
	t.Parallel()
	migrations := []migration{
		{version: 1},
		{version: 2},
		{version: 3},
	}
	tests := []struct {
		name    string
		applied []int
		want    int
		wantErr bool
	}{
		{
			name:    "未適用分を返す",
			applied: []int{1},
			want:    2,
		},
		{
			name:    "履歴が飛んでいる",
			applied: []int{1, 3},
			wantErr: true,
		},
		{
			name:    "履歴が移行数を超える",
			applied: []int{1, 2, 3, 4},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pending, err := pendingMigrations(migrations, tt.applied)
			if (err != nil) != tt.wantErr {
				t.Fatalf("pendingMigrations() error = %v, wantErr %t", err, tt.wantErr)
			}
			if !tt.wantErr && len(pending) != tt.want {
				t.Fatalf("pending migration count = %d, want %d", len(pending), tt.want)
			}
		})
	}
}

func TestV1RelationChecksLimitContradictsAndSupersedesToAssertions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store, err := Open(ctx, filepath.Join(t.TempDir(), "relations.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	for _, tt := range []struct {
		name         string
		relationType string
	}{
		{
			name:         "contradicts",
			relationType: "contradicts",
		},
		{
			name:         "supersedes",
			relationType: "supersedes",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.db.ExecContext(ctx,
				"INSERT INTO relations (relation_id, source_kind, source_id, relation_type, target_kind, target_id, created_at) VALUES (?, 'assertion', 'a1', ?, 'concept', 'c1', '2026-01-01T00:00:00Z')",
				tt.relationType, tt.relationType,
			)
			if err == nil {
				t.Fatalf("invalid %s relation was accepted", tt.relationType)
			}
		})
	}
}

func TestOpenEnablesForeignKeysForStoreConnection(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store, err := Open(ctx, filepath.Join(t.TempDir(), "foreign-keys.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if _, err := store.db.ExecContext(ctx, "INSERT INTO evidence (evidence_id, assertion_id, kind, raw_text, observed_at, created_at) VALUES ('e1', 'missing', 'user_code', 'text', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')"); err == nil {
		t.Fatal("foreign-key violation was accepted")
	}
}

func assertSchemaVersion(t *testing.T, db *sql.DB, want int) {
	t.Helper()
	var got int
	if err := db.QueryRowContext(context.Background(), "SELECT version FROM schema_migrations").Scan(&got); err != nil {
		t.Fatalf("read schema version: %v", err)
	}
	if got != want {
		t.Fatalf("schema version = %d, want %d", got, want)
	}
}

func assertSchemaObjects(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, name := range schemaObjects {
		var found string
		if err := db.QueryRowContext(context.Background(), "SELECT name FROM sqlite_schema WHERE type IN ('table', 'virtual table') AND name = ?", name).Scan(&found); err != nil {
			t.Fatalf("schema object %q missing: %v", name, err)
		}
	}
	for _, name := range []string{
		"idx_scopes_key_value",
		"idx_evidence_assertion",
		"idx_aliases_concept",
		"idx_concept_terms_concept",
		"idx_assertion_concepts_concept",
		"idx_assertion_aliases_value",
		"idx_relations_source",
		"idx_relations_target",
		"idx_temporal_window",
	} {
		var found string
		if err := db.QueryRowContext(context.Background(), "SELECT name FROM sqlite_schema WHERE type = 'index' AND name = ?", name).Scan(&found); err != nil {
			t.Fatalf("index %q missing: %v", name, err)
		}
	}
}

func TestRetrievalOperations(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name   string
		verify func(*testing.T, *Store)
	}{
		{
			name: "字句検索の一致",
			verify: func(t *testing.T, store *Store) {
				results, err := store.SearchText(ctx, "channel")
				if err != nil || len(results) != 1 || results[0].ID != "asrt_01" || !reflect.DeepEqual(results[0].MatchedFields, []string{"assertion_text", "concept", "alias"}) || len(results[0].Concepts) != 1 || len(results[0].Scope) != 1 {
					t.Fatalf("SearchText() = %#v, %v", results, err)
				}
			},
		},
		{
			name: "字句検索の一致なし",
			verify: func(t *testing.T, store *Store) {
				results, err := store.SearchText(ctx, `"`)
				if err != nil || len(results) != 0 {
					t.Fatalf("quoted SearchText() = %#v, %v", results, err)
				}
			},
		},
		{
			name: "Concept Alias検索",
			verify: func(t *testing.T, store *Store) {
				result, err := store.SearchConcept(ctx, "chan")
				if err != nil || result.Concept == nil || result.Concept.ID != "cpt_01" || len(result.Results) != 1 {
					t.Fatalf("SearchConcept() = %#v, %v", result, err)
				}
			},
		},
		{
			name: "Concept検索の一致なし",
			verify: func(t *testing.T, store *Store) {
				result, err := store.SearchConcept(ctx, "missing")
				if err != nil || result.Concept != nil || len(result.Results) != 0 {
					t.Fatalf("missing SearchConcept() = %#v, %v", result, err)
				}
			},
		},
		{
			name: "AssertionなしのConcept",
			verify: func(t *testing.T, store *Store) {
				result, err := store.SearchConcept(ctx, "orphan")
				if err != nil || result.Concept == nil || len(result.Results) != 0 {
					t.Fatalf("orphan SearchConcept() = %#v, %v", result, err)
				}
			},
		},
		{
			name: "Assertion取得",
			verify: func(t *testing.T, store *Store) {
				detail, err := store.GetAssertion(ctx, "asrt_01")
				if err != nil || detail.CurrentRevision != 2 || len(detail.Revisions) != 2 || len(detail.Revisions[0].Scope) != 2 || detail.Revisions[0].Temporal == nil || detail.Revisions[0].Temporal.ObservedAt != nil || detail.Revisions[1].Temporal == nil || len(detail.Concepts) != 1 || !reflect.DeepEqual(detail.Concepts[0].Aliases, []string{"chan"}) || len(detail.Aliases) != 1 {
					t.Fatalf("GetAssertion() = %#v, %v", detail, err)
				}
			},
		},
		{
			name: "Assertion取得のnot_found",
			verify: func(t *testing.T, store *Store) {
				if _, err := store.GetAssertion(ctx, "missing"); !errors.Is(err, domain.ErrAssertionNotFound) {
					t.Fatalf("missing GetAssertion() error = %v", err)
				}
			},
		},
		{
			name: "ScopeなしのAssertion",
			verify: func(t *testing.T, store *Store) {
				detail, err := store.GetAssertion(ctx, "asrt_02")
				if err != nil || len(detail.Revisions) != 1 || detail.Revisions[0].Temporal != nil || len(detail.Revisions[0].Scope) != 0 || len(detail.Concepts) != 0 || len(detail.Aliases) != 0 {
					t.Fatalf("scopeなし GetAssertion() = %#v, %v", detail, err)
				}
			},
		},
		{
			name: "Evidence取得",
			verify: func(t *testing.T, store *Store) {
				result, err := store.GetEvidence(ctx, "asrt_01")
				if err != nil || len(result.Evidence) != 2 || result.Evidence[0].ID != "evd_01" {
					t.Fatalf("GetEvidence() = %#v, %v", result, err)
				}
			},
		},
		{
			name: "Evidenceなし",
			verify: func(t *testing.T, store *Store) {
				result, err := store.GetEvidence(ctx, "asrt_02")
				if err != nil || len(result.Evidence) != 0 {
					t.Fatalf("empty GetEvidence() = %#v, %v", result, err)
				}
			},
		},
		{
			name: "Evidence取得のnot_found",
			verify: func(t *testing.T, store *Store) {
				if _, err := store.GetEvidence(ctx, "missing"); !errors.Is(err, domain.ErrAssertionNotFound) {
					t.Fatalf("missing GetEvidence() error = %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := Open(ctx, filepath.Join(t.TempDir(), "retrieval.db"))
			if err != nil {
				t.Fatalf("Open() error = %v", err)
			}
			t.Cleanup(func() { _ = store.Close() })
			seedRetrievalStore(t, store)
			tt.verify(t, store)
		})
	}
}

func seedRetrievalStore(t *testing.T, store *Store) {
	t.Helper()
	statements := []string{
		"INSERT INTO assertions (assertion_id, current_revision, created_at) VALUES ('asrt_01', 2, '2026-08-14T00:00:00Z')",
		"INSERT INTO assertions (assertion_id, current_revision, created_at) VALUES ('asrt_02', 1, '2026-08-14T00:00:00Z')",
		"INSERT INTO assertion_revisions (assertion_id, revision, normalized_text, created_at) VALUES ('asrt_01', 1, 'old channel', '2026-08-14T00:00:00Z')",
		"INSERT INTO assertion_revisions (assertion_id, revision, normalized_text, created_at) VALUES ('asrt_01', 2, 'channel send', '2026-08-14T00:00:00Z')",
		"INSERT INTO assertion_revisions (assertion_id, revision, normalized_text, created_at) VALUES ('asrt_02', 1, 'empty', '2026-08-14T00:00:00Z')",
		"INSERT INTO revision_scopes (assertion_id, revision, scope_key, scope_value) VALUES ('asrt_01', 1, 'language', 'Go')",
		"INSERT INTO revision_scopes (assertion_id, revision, scope_key, scope_value) VALUES ('asrt_01', 1, 'runtime', 'CLI')",
		"INSERT INTO revision_scopes (assertion_id, revision, scope_key, scope_value) VALUES ('asrt_01', 2, 'language', 'Go')",
		"INSERT INTO temporal_metadata (assertion_id, revision) VALUES ('asrt_01', 1)",
		"INSERT INTO temporal_metadata (assertion_id, revision, observed_at) VALUES ('asrt_01', 2, '2026-08-14T00:00:00Z')",
		"INSERT INTO concepts (concept_id, name, created_at) VALUES ('cpt_01', 'channel', '2026-08-14T00:00:00Z')",
		"INSERT INTO concepts (concept_id, name, created_at) VALUES ('cpt_02', 'orphan', '2026-08-14T00:00:00Z')",
		"INSERT INTO concept_terms (term, concept_id, term_kind) VALUES ('channel', 'cpt_01', 'name')",
		"INSERT INTO concept_terms (term, concept_id, term_kind) VALUES ('chan', 'cpt_01', 'alias')",
		"INSERT INTO concept_terms (term, concept_id, term_kind) VALUES ('orphan', 'cpt_02', 'name')",
		"INSERT INTO concept_aliases (alias, concept_id) VALUES ('chan', 'cpt_01')",
		"INSERT INTO assertion_concepts (assertion_id, concept_id) VALUES ('asrt_01', 'cpt_01')",
		"INSERT INTO assertion_aliases (assertion_id, alias_kind, alias_value) VALUES ('asrt_01', 'identifier', 'channel')",
		"INSERT INTO evidence (evidence_id, assertion_id, kind, raw_text, observed_at, created_at) VALUES ('evd_02', 'asrt_01', 'user_code', 'second', '2026-08-15T00:00:00Z', '2026-08-15T00:00:00Z')",
		"INSERT INTO evidence (evidence_id, assertion_id, kind, raw_text, observed_at, created_at) VALUES ('evd_01', 'asrt_01', 'user_code', 'first', '2026-08-14T00:00:00Z', '2026-08-14T00:00:00Z')",
		"INSERT INTO assertion_lexical_index (assertion_id, normalized_text, concept_name, concept_alias, scope_key, scope_value, assertion_alias) VALUES ('asrt_01', 'channel send', 'channel', 'chan', 'language', 'Go', 'channel')",
	}
	for _, statement := range statements {
		if _, err := store.db.ExecContext(context.Background(), statement); err != nil {
			t.Fatalf("seed query %q: %v", statement, err)
		}
	}
}

func TestRetrievalOperationsReturnStorageErrors(t *testing.T) {
	ctx := context.Background()
	database, err := sql.Open(driverName, filepath.Join(t.TempDir(), "closed-retrieval.db"))
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	store := &Store{db: database}
	for _, tt := range []struct {
		name string
		call func() error
	}{
		{
			name: "SearchText",
			call: func() error {
				_, err := store.SearchText(ctx, "channel")

				return err
			},
		},
		{
			name: "SearchConcept",
			call: func() error {
				_, err := store.SearchConcept(ctx, "channel")

				return err
			},
		},
		{
			name: "GetAssertion",
			call: func() error {
				_, err := store.GetAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name: "GetEvidence",
			call: func() error {
				_, err := store.GetEvidence(ctx, "asrt_01")

				return err
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err == nil {
				t.Fatal("error = nil")
			}
		})
	}
}

func TestRetrievalQueryFailures(t *testing.T) {
	ctx := context.Background()
	coverageQuery = func(string) (driver.Rows, error) {
		return nil, errors.New("取得照会失敗")
	}
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	store := &Store{db: database}

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "字句検索",
			call: func() error {
				_, err := store.SearchText(ctx, "channel")

				return err
			},
		},
		{
			name: "AssertionのConcept",
			call: func() error {
				_, err := store.conceptsForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name: "AssertionのScope",
			call: func() error {
				_, err := store.scopeForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name: "Concept検索",
			call: func() error {
				_, err := store.SearchConcept(ctx, "channel")

				return err
			},
		},
		{
			name: "Assertion取得",
			call: func() error {
				_, err := store.GetAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name: "AssertionのConcept詳細",
			call: func() error {
				_, err := store.conceptDetailsForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name: "AssertionのAlias",
			call: func() error {
				_, err := store.aliasesForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name: "Evidence取得",
			call: func() error {
				_, err := store.GetEvidence(ctx, "asrt_01")

				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err == nil {
				t.Fatal("error = nil")
			}
		})
	}
}

func TestHydrateSummaryScopeFailure(t *testing.T) {
	queryCount := 0
	coverageQuery = func(string) (driver.Rows, error) {
		queryCount++
		if queryCount == 1 {
			return &coverageRows{columns: []string{"concept_id", "name"}}, nil
		}

		return nil, errors.New("Scope照会失敗")
	}
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	if err := (&Store{db: database}).hydrateSummary(context.Background(), &domain.AssertionSummary{ID: "asrt_01"}); err == nil {
		t.Fatal("hydrateSummary() error = nil")
	}
}

func TestRetrievalReadAndIterationFailures(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		columns []string
		call    func(*Store) error
	}{
		{
			name:    "字句検索の行読み込み",
			columns: []string{"assertion_id", "normalized_text", "current_revision", "assertion_text_hit", "concept_hit", "alias_hit", "scope_key_hit", "scope_value_hit"},
			call: func(store *Store) error {
				_, err := store.SearchText(ctx, "channel")

				return err
			},
		},
		{
			name:    "Concept検索の行読み込み",
			columns: []string{"concept_id", "name", "assertion_id", "normalized_text", "current_revision"},
			call: func(store *Store) error {
				_, err := store.SearchConcept(ctx, "channel")

				return err
			},
		},
		{
			name:    "Assertion取得の行読み込み",
			columns: []string{"current_revision", "revision", "normalized_text", "temporal_revision", "valid_from", "valid_until", "version_scope", "observed_at", "last_verified", "scope_key", "scope_value"},
			call: func(store *Store) error {
				_, err := store.GetAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name:    "Conceptの行読み込み",
			columns: []string{"concept_id", "name"},
			call: func(store *Store) error {
				_, err := store.conceptsForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name:    "Scopeの行読み込み",
			columns: []string{"scope_key", "scope_value"},
			call: func(store *Store) error {
				_, err := store.scopeForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name:    "Concept詳細の行読み込み",
			columns: []string{"concept_id", "name", "alias"},
			call: func(store *Store) error {
				_, err := store.conceptDetailsForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name:    "Aliasの行読み込み",
			columns: []string{"alias_kind", "alias_value"},
			call: func(store *Store) error {
				_, err := store.aliasesForAssertion(ctx, "asrt_01")

				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coverageQuery = func(string) (driver.Rows, error) {
				return &coverageRows{
					columns: tt.columns,
					values:  [][]driver.Value{{nil}},
				}, nil
			}
			database, err := sql.Open(coverageDriverName, "")
			if err != nil {
				t.Fatalf("sql.Open() error = %v", err)
			}
			t.Cleanup(func() { _ = database.Close() })
			if err := tt.call(&Store{db: database}); err == nil {
				t.Fatal("error = nil")
			}
		})
	}
}

func TestRetrievalRowsCloseFailures(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		columns []string
		call    func(*Store) error
	}{
		{
			name:    "字句検索",
			columns: []string{"assertion_id", "normalized_text", "current_revision", "assertion_text_hit", "concept_hit", "alias_hit", "scope_key_hit", "scope_value_hit"},
			call: func(store *Store) error {
				_, err := store.SearchText(ctx, "channel")

				return err
			},
		},
		{
			name:    "Concept検索",
			columns: []string{"concept_id", "name", "assertion_id", "normalized_text", "current_revision"},
			call: func(store *Store) error {
				_, err := store.SearchConcept(ctx, "channel")

				return err
			},
		},
		{
			name:    "Assertion取得",
			columns: []string{"current_revision", "revision", "normalized_text", "temporal_revision", "valid_from", "valid_until", "version_scope", "observed_at", "last_verified", "scope_key", "scope_value"},
			call: func(store *Store) error {
				_, err := store.GetAssertion(ctx, "asrt_01")

				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coverageQuery = func(string) (driver.Rows, error) {
				return &coverageRows{
					columns: tt.columns,
					err:     errors.New("行の反復失敗"),
				}, nil
			}
			database, err := sql.Open(coverageDriverName, "")
			if err != nil {
				t.Fatalf("sql.Open() error = %v", err)
			}
			t.Cleanup(func() { _ = database.Close() })
			if err := tt.call(&Store{db: database}); err == nil {
				t.Fatal("error = nil")
			}
		})
	}
}

func TestRetrievalIterationFailures(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		columns []string
		call    func(*Store) error
	}{
		{
			name:    "字句検索",
			columns: []string{"assertion_id", "normalized_text", "current_revision", "assertion_text_hit", "concept_hit", "alias_hit", "scope_key_hit", "scope_value_hit"},
			call: func(store *Store) error {
				_, err := store.SearchText(ctx, "channel")

				return err
			},
		},
		{
			name:    "Concept検索",
			columns: []string{"concept_id", "name", "assertion_id", "normalized_text", "current_revision"},
			call: func(store *Store) error {
				_, err := store.SearchConcept(ctx, "channel")

				return err
			},
		},
		{
			name:    "Assertion取得",
			columns: []string{"current_revision", "revision", "normalized_text", "temporal_revision", "valid_from", "valid_until", "version_scope", "observed_at", "last_verified", "scope_key", "scope_value"},
			call: func(store *Store) error {
				_, err := store.GetAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name:    "Concept",
			columns: []string{"concept_id", "name"},
			call: func(store *Store) error {
				_, err := store.conceptsForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name:    "Scope",
			columns: []string{"scope_key", "scope_value"},
			call: func(store *Store) error {
				_, err := store.scopeForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name:    "Concept詳細",
			columns: []string{"concept_id", "name", "alias"},
			call: func(store *Store) error {
				_, err := store.conceptDetailsForAssertion(ctx, "asrt_01")

				return err
			},
		},
		{
			name:    "Alias",
			columns: []string{"alias_kind", "alias_value"},
			call: func(store *Store) error {
				_, err := store.aliasesForAssertion(ctx, "asrt_01")

				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coverageQuery = func(string) (driver.Rows, error) {
				return &coverageRows{
					columns: tt.columns,
					err:     errors.New("行の反復失敗"),
				}, nil
			}
			database, err := sql.Open(coverageDriverName, "")
			if err != nil {
				t.Fatalf("sql.Open() error = %v", err)
			}
			t.Cleanup(func() { _ = database.Close() })
			if err := tt.call(&Store{db: database}); err == nil {
				t.Fatal("error = nil")
			}
		})
	}
}

func TestEvidenceReadAndIterationFailures(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		rows driver.Rows
		call func(*Store) error
	}{
		{
			name: "行読み込み",
			rows: &coverageRows{
				columns: []string{"evidence_id", "kind", "raw_text", "observed_at"},
				values:  [][]driver.Value{{nil}},
			},
			call: func(store *Store) error {
				_, err := store.GetEvidence(ctx, "asrt_01")

				return err
			},
		},
		{
			name: "行反復",
			rows: &coverageRows{
				columns: []string{"evidence_id", "kind", "raw_text", "observed_at"},
				err:     errors.New("行の反復失敗"),
			},
			call: func(store *Store) error {
				_, err := store.GetEvidence(ctx, "asrt_01")

				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryCount := 0
			coverageQuery = func(string) (driver.Rows, error) {
				queryCount++
				if queryCount == 1 {
					return &coverageRows{
						columns: []string{"exists"},
						values:  [][]driver.Value{{int64(1)}},
					}, nil
				}

				return tt.rows, nil
			}
			database, err := sql.Open(coverageDriverName, "")
			if err != nil {
				t.Fatalf("sql.Open() error = %v", err)
			}
			t.Cleanup(func() { _ = database.Close() })
			if err := tt.call(&Store{db: database}); err == nil {
				t.Fatal("error = nil")
			}
		})
	}
}

func TestMatchedFields(t *testing.T) {
	tests := []struct {
		name string
		got  []string
		want []string
	}{
		{
			name: "本文のみ",
			got:  matchedFields(true, false, false, false, false),
			want: []string{"assertion_text"},
		},
		{
			name: "Conceptのみ",
			got:  matchedFields(false, true, false, false, false),
			want: []string{"concept"},
		},
		{
			name: "Scope keyのみ",
			got:  matchedFields(false, false, false, true, false),
			want: []string{"scope_key"},
		},
		{
			name: "Scope valueのみ",
			got:  matchedFields(false, false, false, false, true),
			want: []string{"scope_value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("matchedFields() = %#v, want %#v", tt.got, tt.want)
			}
		})
	}
}

func TestRetrievalNestedFailures(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		call func(*Store) error
		rows []driver.Rows
	}{
		{
			name: "字句検索の関連取得",
			call: func(store *Store) error {
				_, err := store.SearchText(ctx, "channel")

				return err
			},
			rows: []driver.Rows{
				&coverageRows{
					columns: []string{"assertion_id", "normalized_text", "current_revision", "assertion_text_hit", "concept_hit", "alias_hit", "scope_key_hit", "scope_value_hit"},
					values:  [][]driver.Value{{"asrt_01", "channel", int64(1), true, false, false, false, false}},
				},
			},
		},
		{
			name: "Concept検索のScope取得",
			call: func(store *Store) error {
				_, err := store.SearchConcept(ctx, "channel")

				return err
			},
			rows: []driver.Rows{
				&coverageRows{
					columns: []string{"concept_id", "name", "assertion_id", "normalized_text", "current_revision"},
					values:  [][]driver.Value{{"cpt_01", "channel", "asrt_01", "channel", int64(1)}},
				},
			},
		},
		{
			name: "Assertion取得のConcept取得",
			call: func(store *Store) error {
				_, err := store.GetAssertion(ctx, "asrt_01")

				return err
			},
			rows: []driver.Rows{
				&coverageRows{
					columns: []string{"current_revision", "revision", "normalized_text", "temporal_revision", "valid_from", "valid_until", "version_scope", "observed_at", "last_verified", "scope_key", "scope_value"},
					values:  [][]driver.Value{{int64(1), int64(1), "channel", nil, nil, nil, nil, nil, nil, nil, nil}},
				},
			},
		},
		{
			name: "Evidence照会",
			call: func(store *Store) error {
				_, err := store.GetEvidence(ctx, "asrt_01")

				return err
			},
			rows: []driver.Rows{
				&coverageRows{
					columns: []string{"exists"},
					values:  [][]driver.Value{{int64(1)}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryCount := 0
			coverageQuery = func(string) (driver.Rows, error) {
				if queryCount >= len(tt.rows) {
					return nil, errors.New("後続照会失敗")
				}
				rows := tt.rows[queryCount]
				queryCount++

				return rows, nil
			}
			database, err := sql.Open(coverageDriverName, "")
			if err != nil {
				t.Fatalf("sql.Open() error = %v", err)
			}
			t.Cleanup(func() { _ = database.Close() })
			if err := tt.call(&Store{db: database}); err == nil {
				t.Fatal("error = nil")
			}
		})
	}
}

func TestGetAssertionAliasFailure(t *testing.T) {
	queryCount := 0
	coverageQuery = func(string) (driver.Rows, error) {
		queryCount++
		switch queryCount {
		case 1:
			return &coverageRows{
				columns: []string{"current_revision", "revision", "normalized_text", "temporal_revision", "valid_from", "valid_until", "version_scope", "observed_at", "last_verified", "scope_key", "scope_value"},
				values:  [][]driver.Value{{int64(1), int64(1), "channel", nil, nil, nil, nil, nil, nil, nil, nil}},
			}, nil
		case 2:
			return &coverageRows{columns: []string{"concept_id", "name", "alias"}}, nil
		default:
			return nil, errors.New("Alias照会失敗")
		}
	}
	database, err := sql.Open(coverageDriverName, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	if _, err := (&Store{db: database}).GetAssertion(context.Background(), "asrt_01"); err == nil {
		t.Fatal("GetAssertion() error = nil")
	}
}
