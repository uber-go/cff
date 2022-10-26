# Installation

cff is a library and a CLI.
You can install these independently with `go get`/`go install`
(see [Manual setup](#manual-setup))
but we recommend using the following instructions
to set up your project to use it.

## Recommended setup

### One-time setup

**Prerequisites**

- Go 1.18 or newer
- a project with a go.mod file

Most Go projects should take the following steps to set up cff.

1. Inside the project directory,
   create a 'tools.go' if one doesn't already exist.
   This is where you'll specify development-time dependencies.

   ```bash
   cat > tools.go <<EOF
   //go:build tools

   package tools // use your project's package name here
   EOF
   ```

   Make sure you use the same package name as your project directory.

2. Add `import _ "go.uber.org/cff/cmd/cff"` to the tools.go.

   ```bash
   echo 'import _ "go.uber.org/cff/cmd/cff"' >> tools.go
   ```

3. Run `go mod tidy` to pick up the latest version of cff,
   or run `go get go.uber.org/cff@main` to get the current unreleased branch.

   ```bash
   go mod tidy
   ```

4. Install the cff CLI to a bin/ subdirectory of the project.

   ```bash
   GOBIN=$(pwd)/bin go install go.uber.org/cff/cmd/cff
   ```

   Feel free to gitignore this directory.

   ```bash
   echo '/bin' >> .gitignore
   ```

5. Add the following `go:generate` directive to an existing Go file
   in the same directory.

   ```go
   //go:generate bin/cff ./...
   ```

### Setup on new machines

Once a project is already using cff,
new machines that work on the project simply need to install the cff CLI
into the bin/ directory.

```bash
GOBIN=$(pwd)/bin go install go.uber.org/cff/cmd/cff
```

We recommend incorporating this into your project setup instructions or
scripts.

## Manual setup

Alternatively, you can install the cff CLI and library independently:

1. Add the library as a dependency of the project.

   ```bash
   go get go.uber.org/cff
   ```

2. Install the CLI globally.

   ```bash
   go install go.uber.org/cff/cmd/cff
   ```
