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
	envadds := strings.Split(os.Getenv("INPUT_ENV"), ":::")
	additions := make(map[string]string)
	for _, add := range envadds {
		varsplit := strings.SplitN(add, "=", 2)
		if len(varsplit) != 2 {
			log.Printf("failed to split environment variable '%s' into name and value\n", add)
			os.Exit(1)
		}
		additions[varsplit[0]] = additions[varsplit[1]]
	}

	argrepl := strings.Split(os.Getenv("INPUT_REPLACE_ARGS"), ":::")
	replacements := make(map[string]string)
	for _, rep := range argrepl {
		repsplit := strings.SplitN(rep, "=", 2)
		if len(repsplit) != 2 {
			log.Printf("failed to split argument replacement '%s' into old and new\n", rep)
			os.Exit(1)
		}
		replacements[repsplit[0]] = replacements[repsplit[1]]
	}

	executer := compile_commands_json_executer.Executer{
		Appends:     strings.Split(os.Getenv("INPUT_APPEND_ARGS"), ":"),
		Prepends:    strings.Split(os.Getenv("INPUT_PREPEND_ARGS"), ":"),
		RemoveArgs:  strings.Split(os.Getenv("INPUT_REMOVE_ARGS"), ":"),
		Exe:         os.Getenv("INPUT_EXE"),
		AcceptTU:    strings.Split(os.Getenv("INPUT_ACCEPT_TUS"), ":"),
		RejectTU:    strings.Split(os.Getenv("INPUT_REJECT_TUS"), ":"),
		Env:         additions,
		Replace:     replacements,
		Concurrency: concurrency,
		TraceFile:   os.Getenv("INPUT_TRACE_FILE"),
	}

	err = executer.Run(os.Getenv("INPUT_BUILD_PATH"))
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
}
