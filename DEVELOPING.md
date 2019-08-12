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

To support the possibility of online analysis in the future (vs static
analysis) we want using CFF to feel like using a regular Go library. Users
should find it conceivable, based on the APIs and their behavior, that the
functionality of CFF is completely runtime. The presence of code generation
should be considered an implementation detail by users; it should be hidden
away until absolutely needed. To satisfy this,

- The code generator must only make use of information that could be present
  at runtime. This means that type information may be used, but not variable
  names. Exceptions may be made here for better UX as long as they don't have
  meaningful behavioral impact. For example variable names may be used to
  affect telemetry, but not graph resolution.
- There must be no shared information between Flows unless provided by a
  shared object explicitly. For example, Flows cannot share a global task
  scheduler; instead one must be injected at the `cff.Flow` call site.
  Similarly, there must be no shared global configuration between Flows, but
  configuration may be injected at the call site.

# Testing

For integration testing, test data is placed inside the internal/tests folder.
Each directory in that folder is treated as its own Go package. Write sample
flows in these directories, generate code from those flows, and check in the
generated code. Write tests against this generated code.

The golden_test will verify that all generated code in internal/tests is up to
date.
