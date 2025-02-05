// Copyright 2019 CUE Authors
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

// genversion generates a TypeScript module that contains an exported string
// constant that is the version of the cuelang.org/go module in use.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	// imported for side effect of module being available in cache
	_ "cuelang.org/go/pkg"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// Generate new tour files
	var cueDir bytes.Buffer
	cmd := exec.Command("go", "list", "-m", "-f={{.Version}}", "cuelang.org/go")
	cmd.Stdout = &cueDir
	if err := cmd.Run(); err != nil {
		log.Fatal(fmt.Errorf("failed to run %v; %w", strings.Join(cmd.Args, " "), err))
	}
	out := fmt.Sprintf("export const CUEVersion = \"%v\";\n", strings.TrimSpace(cueDir.String()))
	if err := os.WriteFile("gen_cuelang_org_go_version.ts", []byte(out), 0666); err != nil {
		log.Fatal(fmt.Errorf("failed to write generated version file: %v", err))
	}
}
