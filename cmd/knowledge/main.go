// knowledgeはKnowledge CLIの公開入出力境界を提供する。
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var (
	processArguments           = os.Args
	processStdout    io.Writer = os.Stdout
	processStderr    io.Writer = os.Stderr
	exitProcess                = os.Exit
)

func main() {
	exitProcess(run(processArguments[1:], processStderr))
}

func run(arguments []string, stderr io.Writer) int {
	if _, err := parseCommand(arguments); err.code != "" {
		writeError(stderr, err)

		return exitCode(err.code)
	}

	writeError(stderr, cliError{
		code:    internalError,
		message: "操作の実行処理はまだ利用できません",
	})

	return exitCode(internalError)
}

type successResponse struct {
	OK   bool `json:"ok"`
	Data any  `json:"data"`
}

type errorResponse struct {
	OK    bool      `json:"ok"`
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    errorCode `json:"code"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
}

func writeSuccess(writer io.Writer, data any) {
	writeJSON(writer, successResponse{
		OK:   true,
		Data: data,
	})
}

func writeError(writer io.Writer, err cliError) {
	writeJSON(writer, errorResponse{
		OK: false,
		Error: errorBody{
			Code:    err.code,
			Message: err.message,
			Field:   err.field,
		},
	})
}

func writeJSON(writer io.Writer, value any) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintln(writer, string(encoded))
}

type errorCode string

const (
	validationError errorCode = "validation_error"
	notFoundError   errorCode = "not_found"
	conflictError   errorCode = "conflict"
	storageError    errorCode = "storage_error"
	internalError   errorCode = "internal_error"
)

func exitCode(code errorCode) int {
	switch code {
	case validationError:
		return 2
	case notFoundError:
		return 3
	case conflictError:
		return 4
	case storageError, internalError:
		return 1
	default:
		return 1
	}
}
