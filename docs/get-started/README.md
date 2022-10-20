# Get Started

In this tutorial you will:

- [play with your first cff flow](flow.md)

First, set up a project.

1. Start a new project

   ```bash
   mkdir hellocff
   cd hellocff
   git init
   go mod init example.com/hellocff
   ```

2. Install the cff CLI inside a bin/ subdirectory
   and make sure it doesn't get checked into version-control.

   ```bash
   GOBIN=$(pwd)/bin go install go.uber.org/cff/cmd/cff@latest
   echo '/bin' >> .gitignore
   ```

3. Create a new .go file with a `go:generate` directive to run cff.

   ```bash
   cat > gen.go <<EOF
   package main

   //go:generate bin/cff ./...
   EOF
   ```

Now [write your first flow](flow.md).
