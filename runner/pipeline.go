package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/loader"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/parser"
	"github.com/thenativeweb/esdm/reporter"
	"github.com/thenativeweb/esdm/resolver"
	"github.com/thenativeweb/esdm/rules"
	"github.com/thenativeweb/esdm/schema"
)

// schemasDirectoryDriftRuleID is emitted when a project's
// local schemas/ directory deviates from the embedded
// schema set - an inconsistency the linter refuses to
// continue past, because the lint result would silently
// disagree with what the user sees in their editor.
const schemasDirectoryDriftRuleID = "esdm/system/schemas-directory-drift"

// Run executes the full linter pipeline against the
// directory dir: verify the local schemas/ directory (if
// present) matches the embedded schema set, discover
// sources, parse and schema-validate each file, resolve
// them into a Model, and - if no error-severity
// diagnostics were produced so far - run the rule catalog
// in parallel.
//
// Returned diagnostics are sorted deterministically. The
// error is non-nil only for I/O or infrastructure
// failures (unreadable directory, empty directory,
// broken schema compilation); content-level problems are
// returned as diagnostics.
func Run(ctx context.Context, dir string) ([]diag.Diagnostic, error) {
	diagnostics, _, err := RunWithModel(ctx, dir)
	return diagnostics, err
}

// RunWithModel mirrors Run but additionally returns the
// resolved Model. Consumers that want to introspect the
// model after linting (e.g. the `esdm view` command) use
// this entry point so they do not have to re-run the
// resolver pipeline. The model is non-nil unless an I/O
// or infrastructure error prevented resolution
// altogether.
func RunWithModel(ctx context.Context, dir string) ([]diag.Diagnostic, *model.Model, error) {
	collector := reporter.NewCollector()

	if driftDiagnostic, hasDrift := verifyLocalSchemas(dir); hasDrift {
		collector.Report(driftDiagnostic)
		return collector.All(), nil, nil
	}

	paths, err := loader.Walk(dir)
	if err != nil {
		return nil, nil, err
	}

	parsedFiles := make([]*parser.ParsedFile, 0, len(paths))
	for _, path := range paths {
		parsed, diagnostics, err := parser.Parse(path)
		if err != nil {
			return nil, nil, err
		}

		for _, d := range diagnostics {
			collector.Report(d)
		}

		parsedFiles = append(parsedFiles, parsed)
	}

	resolvedModel, resolverDiagnostics := resolver.Resolve(parsedFiles)
	for _, d := range resolverDiagnostics {
		collector.Report(d)
	}

	if !collector.HasErrors() {
		RunRules(ctx, rules.Catalog(), resolvedModel, collector)
	}

	return collector.All(), resolvedModel, nil
}

// verifyLocalSchemas runs schema.Verify against the
// project's schemas/ directory if present, and translates
// any deviation into a single error diagnostic. The
// second return value indicates whether the linter should
// abort early.
func verifyLocalSchemas(dir string) (diag.Diagnostic, bool) {
	schemasRoot := filepath.Join(dir, "schemas")

	info, err := os.Stat(schemasRoot)
	if err != nil {
		// Missing schemas/ is fine - it is purely an
		// editor convenience. Other stat errors should
		// surface elsewhere (e.g. when Walk runs).
		return diag.Diagnostic{}, false
	}
	if !info.IsDir() {
		return diag.Diagnostic{}, false
	}

	err = schema.Verify(schemasRoot)
	if err != nil {
		var verifyError *schema.VerifyError
		message := err.Error()
		location := diag.Location{File: schemasRoot}
		if errors.As(err, &verifyError) {
			message = fmt.Sprintf("local schemas/ deviates from this binary's schema set: %s. Run `esdm update-schema` to refresh.", verifyError.Error())
			if verifyError.Path != "" {
				location = diag.Location{File: filepath.Join(schemasRoot, filepath.FromSlash(verifyError.Path))}
			}
		}
		return diag.Diagnostic{
			RuleID:   schemasDirectoryDriftRuleID,
			Severity: diag.SeverityError,
			Message:  message,
			Location: location,
		}, true
	}

	return diag.Diagnostic{}, false
}
