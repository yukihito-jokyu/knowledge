package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/yukihito-jokyu/knowledge/internal/persistence/sqlite"
	_ "modernc.org/sqlite"
)

func TestKnowledgeCLIAtProcessBoundary(t *testing.T) {
	fixture := readFixture(t)
	store := defaultStoreConfiguration(t, t.TempDir())
	prepareRetrievalDatabase(t, store.Path, fixture.Seed)
	before, err := os.ReadFile(store.Path)
	if err != nil {
		t.Fatalf("検索前のStoreを読む: %v", err)
	}
	binary := buildCLI(t, false)
	for _, testCase := range fixture.Cases {
		t.Run(testCase.Name, func(t *testing.T) {
			stdout, stderr, err := runCommand(binary, store.Environment, testCase.Arguments)
			if !isExitCode(err, testCase.ExitCode) {
				t.Fatalf("command error = %v, want exit %d", err, testCase.ExitCode)
			}
			assertStdout(t, stdout, testCase.Stdout)
			assertStderr(t, stderr, testCase.Stderr)
		})
	}
	createStore := defaultStoreConfiguration(t, t.TempDir())
	prepareRetrievalDatabase(t, createStore.Path, fixture.Seed)
	for _, testCase := range fixture.CreateCases {
		t.Run(testCase.Name, func(t *testing.T) {
			stdout, stderr, err := runCommand(binary, createStore.Environment, testCase.Arguments)
			if !isExitCode(err, testCase.ExitCode) {
				t.Fatalf("command error = %v, want exit %d", err, testCase.ExitCode)
			}
			assertStderr(t, stderr, testCase.Stderr)
			assertCreateSuccess(t, stdout)
		})
	}
	after, err := os.ReadFile(store.Path)
	if err != nil {
		t.Fatalf("検索後のStoreを読む: %v", err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("検索でStoreが更新されました")
	}
}

func TestKnowledgeCLICreatesDefaultStoreAtProcessBoundary(t *testing.T) {
	fixture := readFixture(t)
	store := defaultStoreConfiguration(t, t.TempDir())
	binary := buildCLI(t, false)
	for _, testCase := range fixture.EmptyStoreCases {
		t.Run(testCase.Name, func(t *testing.T) {
			stdout, stderr, err := runCommand(binary, store.Environment, testCase.Arguments)
			if !isExitCode(err, testCase.ExitCode) {
				t.Fatalf("command error = %v, want exit %d", err, testCase.ExitCode)
			}
			assertStdout(t, stdout, testCase.Stdout)
			assertStderr(t, stderr, testCase.Stderr)
		})
	}
	if _, err := os.Stat(store.Path); err != nil {
		t.Fatalf("既定Storeが作成されない: %v", err)
	}
}

func TestKnowledgeCLIReportsDefaultStoreFailureAtProcessBoundary(t *testing.T) {
	fixture := readFixture(t)
	root := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(root, []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("失敗用設定ルートを作成する: %v", err)
	}
	store := defaultStoreConfiguration(t, root)
	binary := buildCLI(t, false)
	for _, testCase := range fixture.StoreFailureCases {
		t.Run(testCase.Name, func(t *testing.T) {
			stdout, stderr, err := runCommand(binary, store.Environment, testCase.Arguments)
			if !isExitCode(err, testCase.ExitCode) {
				t.Fatalf("command error = %v, want exit %d", err, testCase.ExitCode)
			}
			assertStdout(t, stdout, testCase.Stdout)
			assertStderr(t, stderr, testCase.Stderr)
		})
	}
}

func TestKnowledgeCLICancelsAtProcessBoundary(t *testing.T) {
	fixture := readFixture(t)
	binary := buildCLI(t, true)
	for _, testCase := range fixture.InterruptedCases {
		t.Run(testCase.Name, func(t *testing.T) {
			store := defaultStoreConfiguration(t, t.TempDir())
			if testCase.Stage == "read" || testCase.Stage == "mutation" {
				prepareRetrievalDatabase(t, store.Path, fixture.Seed)
			}
			before := 0
			if testCase.Stage == "mutation" {
				before = assertionCount(t, store.Path)
			}
			stdout, stderr, err := runInterruptedCommand(t, binary, store.Environment, testCase.Arguments, testCase.Stage)
			if !isExitCode(err, testCase.ExitCode) {
				t.Fatalf("command error = %v, want exit %d", err, testCase.ExitCode)
			}
			assertStdout(t, stdout, testCase.Stdout)
			assertStderr(t, stderr, testCase.Stderr)
			if testCase.Stage == "migration" {
				assertMigrationRolledBack(t, store.Path)
			}
			if testCase.Stage == "mutation" && assertionCount(t, store.Path) != before {
				t.Fatal("中断されたmutationがAssertionを残しました")
			}
		})
	}
}

type cliFixture struct {
	Cases             []cliFixtureCase `json:"cases"`
	CreateCases       []cliFixtureCase `json:"create_cases"`
	EmptyStoreCases   []cliFixtureCase `json:"empty_store_cases"`
	StoreFailureCases []cliFixtureCase `json:"store_failure_cases"`
	InterruptedCases  []cliFixtureCase `json:"interrupted_cases"`
	Seed              string           `json:"seed"`
}

type cliFixtureCase struct {
	Name      string         `json:"name"`
	Arguments []string       `json:"arguments"`
	ExitCode  int            `json:"exit_code"`
	Stdout    any            `json:"stdout"`
	Stderr    map[string]any `json:"stderr"`
	Stage     string         `json:"stage"`
}

func readFixture(t *testing.T) cliFixture {
	t.Helper()
	source, err := os.ReadFile(filepath.Join("..", "..", "testdata", "fixtures", "cli-boundary", "cases.json"))
	if err != nil {
		t.Fatalf("fixtureを読む: %v", err)
	}
	var fixture cliFixture
	if err := json.Unmarshal(source, &fixture); err != nil {
		t.Fatalf("fixtureを復号する: %v", err)
	}
	if len(fixture.Cases) == 0 {
		t.Fatal("fixtureにケースがありません")
	}
	if fixture.Seed == "" {
		t.Fatal("fixtureにseedがありません")
	}
	if len(fixture.CreateCases) == 0 || len(fixture.EmptyStoreCases) == 0 || len(fixture.StoreFailureCases) == 0 || len(fixture.InterruptedCases) == 0 {
		t.Fatal("既定Storeのfixtureがありません")
	}

	return fixture
}

func assertCreateSuccess(t *testing.T, source string) {
	t.Helper()
	var response struct {
		OK   bool `json:"ok"`
		Data struct {
			AssertionID string   `json:"assertion_id"`
			Revision    int      `json:"revision"`
			EvidenceIDs []string `json:"evidence_ids"`
			Concepts    []struct {
				ID   string `json:"concept_id"`
				Name string `json:"name"`
			} `json:"concepts"`
			RelationIDs []string `json:"relation_ids"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(source), &response); err != nil {
		t.Fatalf("create stdout JSONを復号する: %v", err)
	}
	if !response.OK || response.Data.AssertionID == "" || response.Data.Revision != 1 || len(response.Data.EvidenceIDs) == 0 {
		t.Fatalf("create response = %#v", response)
	}
}

func assertionCount(t *testing.T, databasePath string) int {
	t.Helper()
	database, err := sql.Open("sqlite", databasePath)
	if err != nil {
		t.Fatalf("Assertion確認用Storeを開く: %v", err)
	}
	defer database.Close()
	var count int
	if err := database.QueryRowContext(context.Background(), "SELECT count(*) FROM assertions").Scan(&count); err != nil {
		t.Fatalf("Assertion数を読む: %v", err)
	}

	return count
}

func buildCLI(t *testing.T, integrationTest bool) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "knowledge")
	arguments := []string{
		"build",
		"-o",
		binary,
	}
	if integrationTest {
		arguments = append(arguments, "-tags", "integrationtest")
	}
	arguments = append(arguments, "../../cmd/knowledge")
	build := exec.CommandContext(context.Background(), "go", arguments...)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v: %s", err, output)
	}

	return binary
}

func runInterruptedCommand(t *testing.T, binary string, environment []string, arguments []string, stage string) (string, string, error) {
	t.Helper()
	readyFile := filepath.Join(t.TempDir(), "ready")
	command := exec.CommandContext(context.Background(), binary, arguments...)
	command.Env = append(environment,
		"KNOWLEDGE_TEST_INTEGRATION_GATE_STAGE="+stage,
		"KNOWLEDGE_TEST_INTEGRATION_GATE_READY="+readyFile,
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		return stdout.String(), stderr.String(), err
	}
	if err := waitForIntegrationGate(readyFile, stage); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()

		return stdout.String(), stderr.String(), err
	}
	if err := command.Process.Signal(os.Interrupt); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()

		return stdout.String(), stderr.String(), err
	}
	err := command.Wait()

	return stdout.String(), stderr.String(), err
}

func waitForIntegrationGate(path string, stage string) error {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		contents, err := os.ReadFile(path)
		if err == nil {
			if string(contents) != stage {
				return fmt.Errorf("integration gate stage = %q, want %q", contents, stage)
			}

			return nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("integration gateを読む: %w", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	return fmt.Errorf("integration gateの待機がタイムアウトしました")
}

func assertMigrationRolledBack(t *testing.T, databasePath string) {
	t.Helper()
	database, err := sql.Open("sqlite", databasePath)
	if err != nil {
		t.Fatalf("migration確認用Storeを開く: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	var count int
	if err := database.QueryRowContext(context.Background(), "SELECT count(*) FROM sqlite_schema WHERE name = 'schema_migrations'").Scan(&count); err != nil {
		t.Fatalf("migration rollbackを確認する: %v", err)
	}
	if count != 0 {
		t.Fatalf("中断されたmigrationが残りました: %d", count)
	}
}

func runCommand(binary string, environment []string, arguments []string) (string, string, error) {
	command := exec.CommandContext(context.Background(), binary, arguments...)
	command.Env = environment
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()

	return stdout.String(), stderr.String(), err
}

func assertStdout(t *testing.T, source string, want any) {
	t.Helper()
	if text, ok := want.(string); ok {
		if source != text {
			t.Fatalf("stdout = %q, want %q", source, text)
		}

		return
	}
	assertSingleJSON(t, source, want, "stdout")
}

func assertStderr(t *testing.T, source string, want map[string]any) {
	t.Helper()
	if want == nil {
		if source != "" {
			t.Fatalf("stderr = %q, want empty", source)
		}

		return
	}
	assertSingleJSON(t, source, want, "stderr")
}

func assertSingleJSON(t *testing.T, source string, want any, stream string) {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewBufferString(source))
	var got any
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("%s JSONを復号する: %v", stream, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		t.Fatalf("%s JSONが複数あります: %v", stream, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s JSON = %#v, want %#v", stream, got, want)
	}
}

func prepareRetrievalDatabase(t *testing.T, databasePath, seed string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(databasePath), 0o700); err != nil {
		t.Fatalf("seed DBの親ディレクトリを作成する: %v", err)
	}
	store, err := sqlite.Open(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("SQLiteを開く: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("SQLiteを閉じる: %v", err)
	}
	database, err := sql.Open("sqlite", databasePath)
	if err != nil {
		t.Fatalf("seed DBを開く: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	source, err := os.ReadFile(filepath.Join("..", "..", "testdata", "fixtures", "cli-boundary", seed))
	if err != nil {
		t.Fatalf("seedを読む: %v", err)
	}
	for _, statement := range strings.Split(string(source), ";") {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}
		if _, err := database.ExecContext(context.Background(), statement); err != nil {
			t.Fatalf("seed DB: %v", err)
		}
	}
}

type defaultStore struct {
	Environment []string
	Path        string
}

func defaultStoreConfiguration(t *testing.T, root string) defaultStore {
	t.Helper()
	switch runtime.GOOS {
	case "darwin":
		return defaultStore{
			Environment: environmentWith(map[string]string{"HOME": root}),
			Path:        filepath.Join(root, "Library", "Application Support", "knowledge", "knowledge.db"),
		}
	case "windows":
		return defaultStore{
			Environment: environmentWith(map[string]string{"APPDATA": root}),
			Path:        filepath.Join(root, "knowledge", "knowledge.db"),
		}
	default:
		return defaultStore{
			Environment: environmentWith(map[string]string{"XDG_CONFIG_HOME": root}),
			Path:        filepath.Join(root, "knowledge", "knowledge.db"),
		}
	}
}

func environmentWith(overrides map[string]string) []string {
	environment := make([]string, 0, len(os.Environ())+len(overrides))
	for _, entry := range os.Environ() {
		key, _, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		if _, overridden := overrides[key]; !overridden {
			environment = append(environment, entry)
		}
	}
	for key, value := range overrides {
		environment = append(environment, key+"="+value)
	}

	return environment
}

func isExitCode(err error, want int) bool {
	if want == 0 {
		return err == nil
	}
	var exitError *exec.ExitError

	return errors.As(err, &exitError) && exitError.ExitCode() == want
}
