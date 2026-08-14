package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseCommandAcceptsAllOperations(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
	}{
		{
			name: "search-text",
			arguments: []string{
				"search-text",
				"--query",
				"channel",
			},
		},
		{
			name: "search-concept",
			arguments: []string{
				"search-concept",
				"--concept",
				"channel",
			},
		},
		{
			name: "search-related",
			arguments: []string{
				"search-related",
				"--seed-kind",
				"assertion",
				"--seed-id",
				"asrt_01",
				"--relation-type",
				"causes",
			},
		},
		{
			name: "get",
			arguments: []string{
				"get",
				"--assertion-id",
				"asrt_01",
			},
		},
		{
			name: "get-evidence",
			arguments: []string{
				"get-evidence",
				"--assertion-id",
				"asrt_01",
			},
		},
		{
			name: "search-contradictions",
			arguments: []string{
				"search-contradictions",
				"--concept",
				"channel",
			},
		},
		{
			name: "search-temporal",
			arguments: []string{
				"search-temporal",
				"--scope-key",
				"language",
				"--scope-value",
				"Go",
			},
		},
		{
			name: "create",
			arguments: []string{
				"create",
				"--normalized-text",
				"channel",
				"--scope-key",
				"language",
				"--scope-value",
				"Go",
				"--concept",
				"channel",
				"--concept-alias",
				"Go channel",
				"--alias-kind",
				"identifier",
				"--alias-value",
				"chan",
				"--relation-type",
				"causes",
				"--relation-target-kind",
				"assertion",
				"--relation-target-id",
				"asrt_00",
				"--evidence-kind",
				"user_code",
				"--evidence-text",
				"evidence",
				"--evidence-observed-at",
				"2026-08-14T00:00:00Z",
				"--valid-from",
				"2026-08-14T00:00:00Z",
			},
		},
		{
			name: "attach-evidence",
			arguments: []string{
				"attach-evidence",
				"--assertion-id",
				"asrt_01",
				"--evidence-kind",
				"user_code",
				"--evidence-text",
				"evidence",
				"--evidence-observed-at",
				"2026-08-14T00:00:00Z",
			},
		},
		{
			name: "revise",
			arguments: []string{
				"revise",
				"--assertion-id",
				"asrt_01",
				"--normalized-text",
				"channel",
			},
		},
		{
			name: "supersede",
			arguments: []string{
				"supersede",
				"--superseded-assertion-id",
				"asrt_01",
				"--replacement-assertion-id",
				"asrt_02",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseCommand(tt.arguments)
			if err.code != "" {
				t.Fatalf("parseCommand() error = %v", err)
			}
			if parsed.operation != tt.arguments[0] {
				t.Fatalf("operation = %q, want %q", parsed.operation, tt.arguments[0])
			}
		})
	}
}

func TestParseCommandRejectsInvalidInput(t *testing.T) {
	tests := []validationCase{
		invalid("操作なし"), invalid("未知の操作", "unknown"), invalid("optionではない値", "get", "asrt_01"), invalid("option名が空", "get", "--"), invalid("値がない", "get", "--assertion-id"), invalid("次のoptionを値にできない", "get", "--assertion-id", "--other"), invalid("この操作では使えないoption", "get", "--query", "text"), invalid("単一値optionの重複", "get", "--assertion-id", "one", "--assertion-id", "two"), invalid("必須optionなし", "get"), invalid("空文字列", "get", "--assertion-id", " "), invalid("Scope group不足", "revise", "--assertion-id", "one", "--normalized-text", "text", "--scope-key", "language"), invalid("Scope keyなし", "revise", "--assertion-id", "one", "--normalized-text", "text", "--scope-value", "Go"), invalid("Scope key重複", "revise", "--assertion-id", "one", "--normalized-text", "text", "--scope-key", "language", "--scope-value", "Go", "--scope-key", "language", "--scope-value", "Go"), invalid("時刻形式不正", "revise", "--assertion-id", "one", "--normalized-text", "text", "--observed-at", "today"), invalid("UTCではない時刻", "revise", "--assertion-id", "one", "--normalized-text", "text", "--observed-at", "2026-08-14T00:00:00+09:00"), invalid("時点の順序不正", "revise", "--assertion-id", "one", "--normalized-text", "text", "--valid-from", "2026-08-15T00:00:00Z", "--valid-until", "2026-08-14T00:00:00Z"), invalid("検索起点種別不正", "search-related", "--seed-kind", "evidence", "--seed-id", "one"), invalid("Relation種別不正", "search-related", "--seed-kind", "concept", "--seed-id", "one", "--relation-type", "invalid"), invalid("矛盾検索のselector不足", "search-contradictions"), invalid("矛盾検索のselector重複", "search-contradictions", "--assertion-id", "one", "--concept", "channel"), invalid("時点検索のselector不足", "search-temporal"), invalid("Concept Aliasが先", "create", "--normalized-text", "text", "--concept-alias", "channel", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Concept重複", "create", "--normalized-text", "text", "--concept", "channel", "--concept", "channel", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Concept Alias重複", "create", "--normalized-text", "text", "--concept", "one", "--concept-alias", "shared", "--concept", "two", "--concept-alias", "shared", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Alias group不足", "create", "--normalized-text", "text", "--alias-kind", "identifier", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Alias種別不正", "create", "--normalized-text", "text", "--alias-kind", "invalid", "--alias-value", "name", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Alias重複", "create", "--normalized-text", "text", "--alias-kind", "identifier", "--alias-value", "name", "--alias-kind", "identifier", "--alias-value", "name", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Relation group不足", "create", "--normalized-text", "text", "--relation-type", "causes", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Relation種別不正", "create", "--normalized-text", "text", "--relation-type", "supersedes", "--relation-target-kind", "assertion", "--relation-target-id", "id", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Relation対象種別不正", "create", "--normalized-text", "text", "--relation-type", "causes", "--relation-target-kind", "invalid", "--relation-target-id", "id", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("contradicts対象不正", "create", "--normalized-text", "text", "--relation-type", "contradicts", "--relation-target-kind", "concept", "--relation-target-id", "id", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Evidence group不足", "create", "--normalized-text", "text"), invalid("Evidence group順序不正", "attach-evidence", "--assertion-id", "one", "--evidence-kind", "user_code", "--evidence-observed-at", "2026-08-14T00:00:00Z", "--evidence-text", "evidence"), invalid("Evidence種別不正", "create", "--normalized-text", "text", "--evidence-kind", "invalid", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalid("Evidence重複", "create", "--normalized-text", "text", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z", "--evidence-kind", "user_code", "--evidence-text", "evidence", "--evidence-observed-at", "2026-08-14T00:00:00Z"), invalidWithCode("同じAssertionの置換", conflictError, "supersede", "--superseded-assertion-id", "one", "--replacement-assertion-id", "one"),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseCommand(tt.arguments)
			want := tt.code
			if want == "" {
				want = validationError
			}
			if err.code != want {
				t.Fatalf("parseCommand() code = %q, want %q", err.code, want)
			}
		})
	}
}

type validationCase struct {
	name      string
	arguments []string
	code      errorCode
}

func invalid(name string, arguments ...string) validationCase {
	return invalidWithCode(name, validationError, arguments...)
}

func invalidWithCode(name string, code errorCode, arguments ...string) validationCase {
	return validationCase{
		name:      name,
		arguments: arguments,
		code:      code,
	}
}

func TestResponsesAndExitCodes(t *testing.T) {
	if got := (cliError{
		message: "失敗",
	}).Error(); got != "失敗" {
		t.Fatalf("Error() = %q", got)
	}

	var success bytes.Buffer
	writeSuccess(&success, map[string]string{"result": "ok"})
	assertJSON(t, success.Bytes(), map[string]any{
		"ok": true,
		"data": map[string]any{
			"result": "ok",
		},
	})

	var failure bytes.Buffer
	writeError(&failure, cliError{
		code:    notFoundError,
		message: "not found",
		field:   "assertion-id",
	})
	assertJSON(t, failure.Bytes(), map[string]any{
		"ok": false,
		"error": map[string]any{
			"code":    "not_found",
			"message": "not found",
			"field":   "assertion-id",
		},
	})

	tests := []struct {
		code errorCode
		want int
	}{
		{
			validationError,
			2,
		},
		{
			notFoundError,
			3,
		},
		{
			conflictError,
			4,
		},
		{
			storageError,
			1,
		},
		{
			internalError,
			1,
		},
		{
			"unknown",
			1,
		},
	}
	for _, tt := range tests {
		if got := exitCode(tt.code); got != tt.want {
			t.Fatalf("exitCode(%q) = %d, want %d", tt.code, got, tt.want)
		}
	}

	var failed bytes.Buffer
	writeJSON(&failed, make(chan int))
	if failed.Len() != 0 {
		t.Fatalf("failed JSON output = %q, want empty", failed.String())
	}
}

func TestMainUsesProcessBoundary(t *testing.T) {
	originalArguments := processArguments
	originalStdout := processStdout
	originalStderr := processStderr
	originalExit := exitProcess
	t.Cleanup(func() {
		processArguments = originalArguments
		processStdout = originalStdout
		processStderr = originalStderr
		exitProcess = originalExit
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := -1
	processArguments = []string{
		"knowledge",
		"get",
	}
	processStdout = &stdout
	processStderr = &stderr
	exitProcess = func(code int) { exitCode = code }

	main()

	if exitCode != 2 || stdout.Len() != 0 || stderr.Len() == 0 {
		t.Fatalf("main() exit/stdout/stderr = %d/%q/%q", exitCode, stdout.String(), stderr.String())
	}
}

func TestRunWritesOnlyOneStream(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		wantCode  int
		wantError errorCode
	}{
		{
			name: "入力不正",
			arguments: []string{
				"get",
			},
			wantCode:  2,
			wantError: validationError,
		},
		{
			name: "未実装操作",
			arguments: []string{
				"get",
				"--assertion-id",
				"asrt_01",
			},
			wantCode:  1,
			wantError: internalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			if got := run(tt.arguments, &stderr); got != tt.wantCode {
				t.Fatalf("run() = %d, want %d", got, tt.wantCode)
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			var response errorResponse
			if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
				t.Fatalf("decode stderr: %v", err)
			}
			if response.Error.Code != tt.wantError {
				t.Fatalf("error code = %q, want %q", response.Error.Code, tt.wantError)
			}
		})
	}
}

func assertJSON(t *testing.T, source []byte, want map[string]any) {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal(source, &got); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("JSON = %#v, want %#v", got, want)
	}
}
