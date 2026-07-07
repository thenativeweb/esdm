package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thenativeweb/esdm/refgen"
)

const snippetBaseDir = "documentation/snippets"

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "refgen:", err)
		os.Exit(1)
	}
}

func run() error {
	snippets, err := refgen.Snippets()
	if err != nil {
		return err
	}
	for _, key := range refgen.SortedPaths(snippets) {
		path := filepath.Join(snippetBaseDir, filepath.FromSlash(key))
		err := os.MkdirAll(filepath.Dir(path), 0o755)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, snippets[key], 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}
