# Go Chapter Scaffold Design

## Goal

Create a lightweight Go project named `agent-plan` so each textbook chapter has a dedicated place for exercises.

## Structure

- `go.mod` defines the module as `agent-plan`.
- `cmd/chapter` provides a tiny CLI for listing chapters and locating one chapter by number.
- `internal/shared` contains common chapter and exercise metadata types.
- `internal/chapter01` through `internal/chapter14` contain chapter-specific metadata, a placeholder `Exercises()` function, and a small compile-time test.
- `Makefile` provides repeatable commands for formatting, testing, listing chapters, and opening one chapter from the CLI.

## Scope

The scaffold intentionally avoids implementing chapter homework. Each chapter package is ready for incremental work without making the initial test suite fail.
