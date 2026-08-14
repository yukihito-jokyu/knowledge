// knowledgeはKnowledge CLIの公開入出力境界を提供する。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
)

var (
	processArguments              = os.Args
	processStdout       io.Writer = os.Stdout
	processStderr       io.Writer = os.Stderr
	exitProcess                   = os.Exit
	parseCLICommand               = parseCommand
	newInterruptContext           = newOSInterruptContext
)

func main() {
	ctx, stop := newInterruptContext()
	code := run(ctx, processArguments[1:], processStderr)
	stop()
	exitProcess(code)
}

func newOSInterruptContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt)
}

func run(ctx context.Context, arguments []string, stderr io.Writer) int {
	return runWithExecutor(ctx, arguments, stderr, executeRetrieval)
}

type retrievalExecutor func(context.Context, command) (any, cliError, bool)

func runWithExecutor(ctx context.Context, arguments []string, stderr io.Writer, execute retrievalExecutor) int {
	if ctx.Err() != nil {
		return interruptedExitCode
	}
	parsed, err := parseCLICommand(arguments)
	if ctx.Err() != nil {
		return interruptedExitCode
	}
	if err.code != "" {
		writeError(stderr, err)

		return exitCode(err.code)
	}
	if data, executionError, handled := execute(ctx, parsed); handled {
		if ctx.Err() != nil {
			return interruptedExitCode
		}
		if executionError.code != "" {
			writeError(stderr, executionError)

			return exitCode(executionError.code)
		}
		writeSuccess(processStdout, data)

		return 0
	}
	if ctx.Err() != nil {
		return interruptedExitCode
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

const interruptedExitCode = 130

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
