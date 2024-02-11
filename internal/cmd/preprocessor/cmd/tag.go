// Copyright 2023 The CUE Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
)

const (
	// tagNorun is the tag used in a txtar-based directive like code or script
	// to indicate that that node should not be run. For an upload directive, it
	// means the commands in the step are not included in the multi-step script that
	// is compiled, but the commands are included in the output. Thus the rendered
	// output is just the commands, but not any output they would otherwise have
	// generated.
	tagNorun = "norun"

	// tagNoFormat is the tag used in a txtar-based directive to indicate that a
	// file in the archive should not be formatted. A tagNoFormat tag requires an
	// argument, the filepath of the file not to format.
	tagNoFmt = "nofmt"

	// tagCodeTab identifies the tag key used to pass options to the code-tab emitted
	// for a file in a txtar-based directive like code or upload. A tagCodeTab
	// tag requires an argument, the filepath of a file in the files in the archive.
	tagCodeTab = "codetab"

	// tagLocation identifies the location key used to define the location of
	// files in a txtar archive. If specified, the tagLocation requires as many
	// unquoted arguments as there are files. e.g.
	//
	//     #location top-left top-right bottom
	//
	tagLocation = "location"

	tagEllipsis = "ellipsis"
)

// findTag searches for the first #$key (or #$key($arg) if arg is non empty)
// tag line in src. Tags are # prefixed lines where the # at the beginning of
// the line must be followed by a non-space character. args contains the
// contents of the quote-aware args that follow the tag name. present indicates
// whether the tag identified by key was present or not. err will be non-nil if
// there were errors in parsing the arguments to a tag.
//
// TODO: work out whether we want to handle comments in tag lines (which are
// themselves comments already).
//
// TODO: add an explicit test for when arg != ""
func findTag(src []byte, key, arg string) (args []string, present bool, err error) {
	prefix := "#" + key
	if arg != "" {
		prefix += "(" + arg + ")"
	}
	sc := bufio.NewScanner(bytes.NewReader(src))
	lineNo := 1
	for sc.Scan() {
		line := sc.Bytes()
		if after, found := bytes.CutPrefix(bytes.TrimSpace(line), []byte(prefix)); found {
			args, err := parseLineArgs(string(after))
			if err != nil {
				err = fmt.Errorf("%d %w", lineNo, err)
			}
			return args, true, err
		}
		lineNo++
	}
	return nil, false, nil
}
