# Go Modules & Packages Wiki

---

## Table of Contents

1. [Go Module](#go-module)
2. [Go Package](#go-package)
3. [Module vs Package](#module-vs-package)
4. [Examples](#examples)

---

## Go Module

### What is a Module?

A module is a **collection of related Go packages** that are versioned and distributed together. It is defined by a `go.mod` file at its root.

### Why Modules Exist

Before modules (pre-Go 1.11), all Go code lived in a single `$GOPATH` directory with no versioning. Modules solved:

- **Dependency versioning** — lock specific versions of libraries
- **Reproducible builds** — everyone gets the same dependency versions
- **No more GOPATH restriction** — code can live anywhere on disk
- **Dependency isolation** — different projects can use different versions of the same library

### Key Files

| File | Purpose |
|------|---------|
| `go.mod` | Declares module name, Go version, and dependencies |
| `go.sum` | Checksums of dependencies (ensures integrity, like a lockfile) |

### `go.mod` Structure

```
module github.com/yourname/project   ← module identity
go 1.21                               ← minimum Go version
require (
    github.com/sirupsen/logrus v1.9.3  ← dependency with version
    github.com/gin-gonic/gin v1.9.1
)
```

### Common Commands

| Command | What it does |
|---------|-------------|
| `go mod init <module-path>` | Create a new module (generates `go.mod`) |
| `go mod tidy` | Add missing deps, remove unused ones |
| `go mod download` | Download all deps to local cache |
| `go mod verify` | Verify deps haven't been tampered with |
| `go mod graph` | Print dependency graph |
| `go mod vendor` | Copy deps into a `vendor/` folder |

### Module Path Conventions

| Path Style | When to Use |
|-----------|-------------|
| `github.com/user/repo` | Public or private GitHub repos |
| `company.com/team/project` | Corporate/internal projects |
| `myproject` | Local-only, never shared |

### Use Cases

#### 1. Building a Library (Shared Package)

```
my-logger/
├── go.mod          ← module github.com/yourname/my-logger
├── logger.go       ← package logger
└── formatter/
    └── format.go   ← package formatter
```

Others import:
```go
import "github.com/yourname/my-logger"
import "github.com/yourname/my-logger/formatter"
```

#### 2. Building an Executable (CLI Tool / Web Server)

```
my-app/
├── go.mod          ← module github.com/yourname/my-app
├── main.go         ← package main (has func main())
└── handlers/
    └── api.go      ← package handlers
```

#### 3. Private Module

```bash
# Tell Go to skip the public proxy for your private code
export GOPRIVATE=github.com/yourcompany/*

# Use SSH for authentication
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

#### 4. Multi-Module Monorepo

```
monorepo/
├── service-a/
│   ├── go.mod      ← module github.com/company/monorepo/service-a
│   └── main.go
├── service-b/
│   ├── go.mod      ← module github.com/company/monorepo/service-b
│   └── main.go
└── shared/
    ├── go.mod      ← module github.com/company/monorepo/shared
    └── utils.go
```

---

## Go Package

### What is a Package?

A package is a **folder of `.go` files** that all share the same `package` declaration. It's the basic unit of code organization in Go.

### Rules

1. **All `.go` files in the same folder must have the same package name**
2. **One folder = one package** (except `_test` suffix for test files)
3. **Exported names start with uppercase** (`Bark` is public, `bark` is private)
4. **Package name should be lowercase**, no underscores, no camelCase

### Two Types of Packages

#### 1. `package main` — Executable

```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("I'm a runnable program!")
}
```

- Must have `func main()`
- `go run main.go` or `go build` produces a binary
- **One per directory**

#### 2. Named Package (e.g., `package puppy`) — Library

```go
// puppy.go
package puppy

func Bark() string {
    return "Woff!!"
}
```

- No `func main()`
- Cannot be run directly
- Imported by other packages

### Package vs Folder Name

The package name does **not** have to match the folder name, but by convention it should:

```
animals/
└── dog/
    └── dog.go      ← package dog  ✅ (matches folder)
```

This is valid but confusing:
```
animals/
└── dog/
    └── dog.go      ← package cat  ⚠️ (works but bad practice)
```

### Importing Packages

```go
import "fmt"                                          // standard library
import "github.com/yourname/my-logger"                // external module
import "github.com/yourname/my-app/handlers"          // sub-package within your module
```

### Exported vs Unexported

```go
package puppy

func Bark() string { ... }     // Exported  — accessible from outside (uppercase)
func bark() string { ... }     // Unexported — only usable inside this package (lowercase)

var Name string                // Exported
var age int                    // Unexported
```

### Use Cases

#### 1. Utility Package

```
myapp/
├── go.mod
├── main.go              ← package main
└── utils/
    └── helpers.go        ← package utils
```

```go
// utils/helpers.go
package utils

func FormatName(first, last string) string {
    return first + " " + last
}
```

```go
// main.go
package main

import (
    "fmt"
    "github.com/yourname/myapp/utils"
)

func main() {
    fmt.Println(utils.FormatName("John", "Doe"))
}
```

#### 2. Multiple Files in One Package

```
puppy/
├── go.mod
├── puppy.go             ← package puppy (Bark function)
└── tricks.go            ← package puppy (Sit, Fetch functions)
```

All functions from both files are available as `puppy.Bark()`, `puppy.Sit()`, etc.
They can also call each other's unexported functions since they're in the same package.

#### 3. Nested Packages

```
zoo/
├── go.mod               ← module github.com/yourname/zoo
├── zoo.go               ← package zoo
├── animals/
│   ├── animals.go       ← package animals
│   ├── dog/
│   │   └── dog.go       ← package dog
│   └── cat/
│       └── cat.go       ← package cat
└── tickets/
    └── tickets.go       ← package tickets
```

```go
import "github.com/yourname/zoo"
import "github.com/yourname/zoo/animals"
import "github.com/yourname/zoo/animals/dog"
import "github.com/yourname/zoo/animals/cat"
import "github.com/yourname/zoo/tickets"
```

#### 4. Test Package

```
puppy/
├── puppy.go             ← package puppy
├── puppy_test.go        ← package puppy      (white-box test, can access unexported)
└── puppy_ext_test.go    ← package puppy_test  (black-box test, only exported)
```

---

## Module vs Package

| | Module | Package |
|---|---|---|
| Defined by | `go.mod` file | `package` keyword in `.go` files |
| Scope | Entire project (can contain many packages) | Single folder of `.go` files |
| Purpose | Versioning, distribution, dependency management | Code organization, encapsulation |
| Identity | Module path (`github.com/user/repo`) | Package name (`package puppy`) |
| Analogy | A book | A chapter in the book |

### Visual Diagram

```
MODULE (go.mod)
└── github.com/yourname/myapp
    │
    ├── PACKAGE: main        (main.go)          ← entry point
    ├── PACKAGE: handlers    (handlers/*.go)     ← HTTP handlers
    ├── PACKAGE: models      (models/*.go)       ← data structures
    └── PACKAGE: utils       (utils/*.go)        ← helper functions
```

---

## Quick Reference

```bash
# Create a module
go mod init github.com/yourname/project

# Add/remove dependencies automatically
go mod tidy

# Run an executable (package main)
go run .
go run main.go

# Build a binary
go build -o myapp .

# Install a dependency
go get github.com/some/library@v1.2.3

# List all dependencies
go list -m all
```
