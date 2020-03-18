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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pseyfert/compile_commands_json_executer/lib"
)

func modSplit(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

func main() {
	concurrency, err := strconv.Atoi(os.Getenv("INPUT_CONCURRENCY"))
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
	envadds := modSplit(os.Getenv("INPUT_ENV"), ":::")
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

	argrepl := modSplit(os.Getenv("INPUT_REPLACE_ARGS"), ":::")
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
		Appends:     modSplit(os.Getenv("INPUT_APPEND_ARGS"), ":"),
		Prepends:    modSplit(os.Getenv("INPUT_PREPEND_ARGS"), ":"),
		RemoveArgs:  modSplit(os.Getenv("INPUT_REMOVE_ARGS"), ":"),
		Exe:         os.Getenv("INPUT_EXE"),
		AcceptTU:    modSplit(os.Getenv("INPUT_ACCEPT_TUS"), ":"),
		RejectTU:    modSplit(os.Getenv("INPUT_REJECT_TUS"), ":"),
		Env:         additions,
		Replace:     replacements,
		Concurrency: concurrency,
		TraceFile:   os.Getenv("INPUT_TRACE_FILE"),
	}

	database := os.Getenv("INPUT_BUILD_PATH")
	if wd, err := os.Getwd(); err == nil {
		log.Printf("Will try to read %s, which might be relative to %s\n", database, wd)
	} else {
		log.Printf("Will try to read %s, without usable working directory: %v\n", database, err)
	}
	err = githubenv()
	if err != nil {
		log.Printf("bodging github working directories failed: %v\ngood luck.", err)
	}
	err = executer.Run(database)
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
}

func githubenv() error {
	var err error
	reponame := os.Getenv("GITHUB_REPOSITORY")
	split := strings.SplitN(reponame, "/", 2)
	if len(split) != 2 {
		return fmt.Errorf("could not parse repository name %s", reponame)
	}

	gooddir := "__w/" + split[1]
	betterworkdir := gooddir + "/" + split[1]
	if err = os.MkdirAll(gooddir, 0755); err != nil {
		return err
	}
	badworkdir := os.Getenv("GITHUB_WORKSPACE")
	if err = os.Symlink(badworkdir, betterworkdir); err != nil {
		return err
	}
	log.Printf("created dir %s and symlink therein %s pointing to %s\n", gooddir, betterworkdir, badworkdir)
	if _, err := os.Stat(gooddir); os.IsNotExist(err) {
		return fmt.Errorf("new dir doesn't exist")
	}
	if _, err := os.Stat(betterworkdir); os.IsNotExist(err) {
		return fmt.Errorf("new symlink doesn't exist")
	}
	if _, err := os.Stat(badworkdir); os.IsNotExist(err) {
		return fmt.Errorf("symlink target doesn't exist")
	}
	if files, err := ioutil.ReadDir(gooddir); err != nil {
		atleastonefile := false
		for _, file := range files {
			log.Printf("in the created dir, there is a file: %s\n", file)
			atleastonefile = true
		}
		if !atleastonefile {
			log.Printf("there appears to be no content in %s\n", gooddir)
		}
	} else {
		return fmt.Errorf("could not read contents of %s", gooddir)
	}
	if files, err := ioutil.ReadDir(betterworkdir); err != nil {
		atleastonefile := false
		for _, file := range files {
			log.Printf("in the symlinked dir, there is a file: %s\n", file)
			atleastonefile = true
		}
		if !atleastonefile {
			log.Printf("there appears to be no content in %s\n", betterworkdir)
		}
	} else {
		return fmt.Errorf("could not read contents of %s", betterworkdir)
	}
	return nil
	// log.Println("Environment is")
	// for _, e := range os.Environ() {
	// 	pair := strings.SplitN(e, "=", 2)
	// 	if len(pair) == 2 {
	// 		log.Printf("%s\t: %s\n", pair[0], pair[1])
	// 	} else {
	// 		log.Printf("Error when parsing %s\n", e)
	// 	}
	// }
}
