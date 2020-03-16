package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pseyfert/compile_commands_json_executer/lib"
)

func main() {
	concurrency, err := strconv.Atoi(os.Getenv("INPUT_CONCURRENCY"))
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
	executer := compile_commands_json_executer.Executer{
		Appends:    strings.Split(os.Getenv("INPUT_APPEND_ARGS"), ":"),
		Prepends:   strings.Split(os.Getenv("INPUT_PREPEND_ARGS"), ":"),
		RemoveArgs: strings.Split(os.Getenv("INPUT_REMOVE_ARGS"), ":"),
		Exe:        os.Getenv("INPUT_EXE"),
		AcceptTU:   strings.Split(os.Getenv("INPUT_ACCEPT_TUS"), ":"),
		RejectTU:   strings.Split(os.Getenv("INPUT_REJECT_TUS"), ":"),
		// Env:         os.Getenv(""),
		// Replace:     os.Getenv(""),
		Concurrency: concurrency,
		TraceFile:   os.Getenv("INPUT_TRACE_FILE"),
	}

	err = executer.Run(os.Getenv("INPUT_BUILD_PATH"))
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
}
