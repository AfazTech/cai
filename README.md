# CAI

**CAI (Codebase AI Interface)** is a lightweight and practical CLI tool for generating AI-friendly snapshots of your codebase.

It creates a structured text representation of your project, including the project tree and source files, making it easy to provide your entire codebase as context to an LLM for code analysis, debugging, refactoring, and development.

## Donate

If CAI helps you, consider supporting the project and its development.

Your support helps keep the project maintained and allows new features and improvements to be developed.

## Table of Contents

* [Overview](#overview)
* [Features](#features)
* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
* [Project Configuration](#project-configuration)
* [Supported Include Groups](#supported-include-groups)
* [Contributing](#contributing)
* [License](#license)

## Overview

CAI is designed to make it easier to share a complete codebase with AI assistants.

Instead of manually copying files and directories, CAI generates a single structured output containing:

* Project structure
* Source file contents
* Project description
* Configurable ignore patterns
* Configurable include patterns
* Optional `.gitignore` support
* Configurable maximum output size

The generated snapshot is optimized for use as context when working with LLMs.

## Features

* Generate a complete codebase snapshot with a single command
* Automatically generate a project directory tree
* Include source files in the generated snapshot
* Configure which files and directories should be ignored
* Configure which files should be included
* Support regular expressions with `regex:` patterns
* Support predefined language and technology groups
* Optional `.gitignore` support
* Configurable maximum output size
* Automatically detect the project associated with the current directory
* Lightweight and fast Go CLI
* Cross-platform binaries
* Easy installation through `install.sh`

## Requirements

CAI is written in Go and requires:

* Go 21 or higher for building from source
* Linux, macOS, or Windows

Pre-built binaries are available through GitHub Releases, so Go is **not required** when installing a released version.

## Installation

### Quick Install

On Linux and macOS, you can install the latest release using:

```bash
curl -fsSL https://raw.githubusercontent.com/AfazTech/cai/main/install.sh | bash
```

The installer automatically detects your operating system and architecture, downloads the appropriate binary from the latest GitHub Release, and installs `cai` into:

```text
/usr/local/bin
```

If elevated permissions are required, the installer uses `sudo` automatically.

After installation, verify that CAI is available:

```bash
cai
```

### Install from Source

Clone the repository:

```bash
git clone https://github.com/AfazTech/cai.git
cd cai
```

Install dependencies:

```bash
go mod tidy
```

Build the binary:

```bash
go build -o cai ./cmd/cai
```

You can then run CAI directly:

```bash
./cai
```

## Usage

CAI supports several commands for managing projects and generating codebase snapshots.

### Initialize a Project

Initialize CAI configuration for the current directory:

```bash
cai init my-project
```

This creates a project configuration that CAI can use when generating snapshots.

### Generate a Snapshot

Generate a snapshot for the current directory:

```bash
cai .
```

You can also generate a snapshot for a specific directory:

```bash
cai /path/to/project
```

### Use a Specific Project Configuration

```bash
cai -project my-project /path/to/project
```

The generated snapshot is saved using the project name:

```text
my-project.txt
```

### Set a Project Description

```bash
cai set-description "A modern Go CLI application"
```

Short alias:

```bash
cai sd "A modern Go CLI application"
```

### Set Maximum Output Size

Set the maximum generated snapshot size in MB:

```bash
cai set-size 50
```

Short alias:

```bash
cai ss 50
```

### List Projects

Display all configured projects:

```bash
cai list-projects
```

Short alias:

```bash
cai lp
```

## Ignore Patterns

You can add files or directories to the ignore list.

For example:

```bash
cai add-ignore node_modules
```

Short alias:

```bash
cai ai node_modules
```

Remove an ignore pattern:

```bash
cai remove-ignore node_modules
```

Short alias:

```bash
cai ri node_modules
```

CAI also supports regular expressions by using the `regex:` prefix.

Example:

```bash
cai add-ignore 'regex:^tmp/.*'
```

## Include Patterns

You can control which files are included in the generated snapshot.

Add an include pattern:

```bash
cai add-include '*.go'
```

Short alias:

```bash
cai adi '*.go'
```

Remove an include pattern:

```bash
cai remove-include '*.go'
```

Short alias:

```bash
cai rmi '*.go'
```

## Include Groups

CAI provides predefined include groups for common programming languages and technologies.

For example, to include PHP-related files:

```bash
cai add-include-group php
```

Short alias:

```bash
cai aig php
```

To remove a group:

```bash
cai remove-include-group php
```

Short alias:

```bash
cai rig php
```

Available groups include:

* `frontend`
* `markup`
* `styles`
* `js`
* `ts`
* `react`
* `vue`
* `svelte`
* `astro`
* `php`
* `python`
* `c`
* `cpp`
* `csharp`
* `java`
* `go`
* `rust`
* `ruby`
* `swift`
* `kotlin`
* `scala`
* `perl`
* `lua`
* `r`
* `shell`
* `sql`
* `docker`
* `yaml`
* `json`
* `toml`
* `ini`
* `make`
* `cmake`

## Project Configuration

CAI stores project configuration inside:

```text
.cai/
```

Each project has its own configuration file:

```text
.cai/projects/<project>/config.json
```

A configuration can define:

* Project name
* Project root path
* Whether to generate a project tree
* Whether to include file contents
* Ignore patterns
* Include patterns
* `.gitignore` support
* Project description
* Maximum output size

Example configuration:

```json
{
  "project": "my-project",
  "rootPath": "/path/to/project",
  "tree": true,
  "files": true,
  "ignore": [
    ".git",
    "node_modules",
    "vendor",
    "dist",
    "build",
    ".env"
  ],
  "include": [
    "*"
  ],
  "gitignore": false,
  "description": "My project",
  "maxSizeMB": 50
}
```

## Output Format

The generated snapshot contains a structured representation of the project.

It includes the project structure:

```text
# PROJECT STRUCTURE

.
├── cmd
│   └── cai
│       └── main.go
└── internal
    └── core
        ├── app.go
        ├── commands.go
        └── config.go
```

And the contents of included files:

```text
@@@FILE: cmd/cai/main.go@@@

package main

import (
    "os"
    "github.com/AfazTech/cai/internal/core"
)
```

This format makes the generated output easy to provide as context to an AI assistant.

## Releases

CAI uses GitHub Actions to automatically build and publish releases.

Creating and pushing a version tag triggers the release workflow:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The release workflow builds binaries for supported platforms and publishes them as GitHub Release assets.

## Contributing

Contributions are always welcome!

To contribute:

1. Fork the repository.
2. Clone your fork.
3. Create a new branch.
4. Make your changes.
5. Test your changes.
6. Commit your changes.
7. Push your branch.
8. Open a Pull Request.

Example:

```bash
git checkout -b feature/my-feature

git add .

git commit -m "Add my feature"

git push origin feature/my-feature
```

Please keep contributions focused and update documentation when introducing new features or changing existing behavior.

## License

CAI is licensed under the MIT License.

See the [LICENSE](LICENSE) file for more information.
