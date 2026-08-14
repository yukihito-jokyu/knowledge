package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
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
	columns []string
	values  [][]driver.Value
	err     error
	index   int
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
	return nil
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
	t.Cleanup(func() { commitMigration = originalCommit })

	tests := []struct {
		name       string
		migrations []migration
		commitErr  error
		wantErr    bool
	}{
		{name: "移行なし"},
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
			err := store.applyMigrationScripts(ctx, tt.migrations)
			if (err != nil) != tt.wantErr {
				t.Fatalf("applyMigrationScripts() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
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

	for _, relationType := range []string{
		"contradicts",
		"supersedes",
	} {
		_, err := store.db.ExecContext(ctx,
			"INSERT INTO relations (relation_id, source_kind, source_id, relation_type, target_kind, target_id, created_at) VALUES (?, 'assertion', 'a1', ?, 'concept', 'c1', '2026-01-01T00:00:00Z')",
			relationType, relationType,
		)
		if err == nil {
			t.Fatalf("invalid %s relation was accepted", relationType)
		}
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
