package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestKnowledgeCLIAtProcessBoundary(t *testing.T) {
	fixture := readFixture(t)
	binary := buildCLI(t)
	for _, testCase := range fixture.Cases {
		t.Run(testCase.Name, func(t *testing.T) {
			stdout, stderr, err := runCommand(binary, testCase.Arguments)
			if !isExitCode(err, testCase.ExitCode) {
				t.Fatalf("command error = %v, want exit %d", err, testCase.ExitCode)
			}
			if stdout != testCase.Stdout {
				t.Fatalf("stdout = %q, want %q", stdout, testCase.Stdout)
			}
			assertSingleJSON(t, stderr, testCase.Stderr)
		})
	}
}

type cliFixture struct {
	Cases []cliFixtureCase `json:"cases"`
}

type cliFixtureCase struct {
	Name      string         `json:"name"`
	Arguments []string       `json:"arguments"`
	ExitCode  int            `json:"exit_code"`
	Stdout    string         `json:"stdout"`
	Stderr    map[string]any `json:"stderr"`
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

	return fixture
}

func buildCLI(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "knowledge")
	build := exec.CommandContext(context.Background(), "go", "build", "-o", binary, "../../cmd/knowledge")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v: %s", err, output)
	}

	return binary
}

func runCommand(binary string, arguments []string) (string, string, error) {
	command := exec.CommandContext(context.Background(), binary, arguments...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()

	return stdout.String(), stderr.String(), err
}

func assertSingleJSON(t *testing.T, source string, want map[string]any) {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewBufferString(source))
	var got map[string]any
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("stderr JSONを復号する: %v", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		t.Fatalf("stderr JSONが複数あります: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stderr JSON = %#v, want %#v", got, want)
	}
}

func isExitCode(err error, want int) bool {
	var exitError *exec.ExitError

	return errors.As(err, &exitError) && exitError.ExitCode() == want
}
