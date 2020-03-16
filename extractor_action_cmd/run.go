/*
 * Copyright (C) 2020 Paul Seyfert <pseyfert.mathphys@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
		if add == "" {
			continue
		}
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
		if rep == "" {
			continue
		}
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
