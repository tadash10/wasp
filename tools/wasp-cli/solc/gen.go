// go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
	"time"
)

// This program generates solc.go. It can be invoked by running
// go generate from tools/wasp-cli

var (
	re              = regexp.MustCompile(`(?m)(?P<version>\d+\.\d+\.\d+\+commit\.\w+)`)
	packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package util

const SolcVersion = "{{ .SOLCVersion }}"

`))
)

func main() {
	f, err := os.Create("util/solc.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	v := os.Getenv("SOLC_VERSION")

	if err := packageTemplate.Execute(f, struct {
		Timestamp   time.Time
		SOLCVersion string
	}{
		Timestamp:   time.Now(),
		SOLCVersion: re.FindString(v),
	}); err != nil {
		panic(err)
	}

	fmt.Printf("Compiler run successful. Artifact(s) can be found in directory \"%s\".\n", filepath.Dir(f.Name()))
}
