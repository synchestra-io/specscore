// Package idea parses and represents SpecScore Idea artifacts.
//
// An Idea is a single markdown file under `spec/ideas/` (or
// `spec/ideas/archived/`) that captures a pre-spec, lintable one-pager:
// problem framing, recommended direction, MVP scope, exclusions, and
// dealbreaker assumptions. See `spec/features/idea/README.md` for the full
// artifact schema.
//
// This package provides parsing utilities — the lint rules that enforce the
// schema live in `pkg/lint`.
package idea
