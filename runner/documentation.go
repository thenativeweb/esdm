// Package runner ties the linter pipeline together:
// loading files, parsing and schema-validating them,
// resolving them into a model, and - if the earlier
// stages produced no errors - running the rule engine
// with parallel rule execution and panic isolation.
//
// The pipeline is deliberately linear: each stage's
// output feeds the next. An error from the parser or
// resolver stage short-circuits the rule engine, because
// running rules against an incomplete or malformed model
// would produce noisy, confusing output.
package runner
