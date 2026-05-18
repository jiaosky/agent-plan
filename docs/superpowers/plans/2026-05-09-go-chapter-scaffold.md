# Go Chapter Scaffold Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go module named `agent-plan` with one package per textbook chapter.

**Architecture:** Keep each chapter isolated under `internal/chapterXX`, with common metadata types in `internal/shared`. A small CLI in `cmd/chapter` lists available chapters and helps jump into a chosen chapter.

**Tech Stack:** Go standard library, `go test`, `go fmt`, `make`.

---

### Task 1: Module and Shared Types

**Files:**
- Create: `go.mod`
- Create: `internal/shared/types.go`

- [x] **Step 1: Define the module**

```go
module agent-plan

go 1.24
```

- [x] **Step 2: Add shared chapter metadata types**

```go
package shared

type Chapter struct {
	Number string
	Title  string
	Source string
}

type Exercise struct {
	Name        string
	Description string
	Done        bool
}

func (c Chapter) IsValid() bool {
	return c.Number != "" && c.Title != "" && c.Source != ""
}
```

### Task 2: Chapter Packages

**Files:**
- Create: `internal/chapter01` through `internal/chapter14`

- [x] **Step 1: Add one metadata file per chapter**

Each package exposes `Chapter()` and `Exercises()`.

- [x] **Step 2: Add one basic test per chapter**

Each package verifies chapter metadata is populated and the exercise placeholder is callable.

### Task 3: CLI and Developer Commands

**Files:**
- Create: `cmd/chapter/main.go`
- Create: `Makefile`

- [x] **Step 1: Add CLI**

`go run ./cmd/chapter` lists all chapters. `go run ./cmd/chapter -chapter 05` prints the source document and package path.

- [x] **Step 2: Add Makefile**

`make test`, `make fmt`, `make list`, and `make run-chapter CH=05` wrap common commands.

### Task 4: Verify

- [x] **Step 1: Run formatting**

```bash
go fmt ./...
```

- [x] **Step 2: Run tests**

```bash
go test ./...
```
