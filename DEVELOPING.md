# Phases

CFF does its work in two phases: compilation and generation. During
compilation, the user-provided code is inspected and compiled into an internal
representation. Any requisite validation is done at this time. During
generation, the previously built internal representation is used to generate
code.

# Design Principles

We tried to follow some standard best practices in designing this tool.

- Do One Thing and Do It Well. (Unix Philosophy)
- Explicit is better than implicit. (The Zen of Python)

Keeping that in mind, the following guiding principles have helped make
decisions.

- No assumptions should be made, nor restrictions imposed, with regards to code
  organization or abstraction design of user-owned business logic. We are only
  orchestrating functions so we should not need anything except function
  pointers.
- No new dependencies should be added if they can be avoided. For example, a
  dependency on the logging library should not be added if the user did not ask
  for logging.
- No new syntax may be introduced. Both, the input and output files must be
  valid Go code that passes type checking.

# Testing

For integration testing, test data is placed inside the internal/tests folder.
Each directory in that folder is treated as its own Go package. Write sample
flows in these directories, generate code from those flows, and check in the
generated code. Write tests against this generated code.

The golden_test will verify that all generated code in internal/tests is up to
date.
