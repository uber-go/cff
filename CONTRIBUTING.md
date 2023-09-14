---
sidebarDepth: 2
search: false
---

# Contributing

Thanks for trying to help make cff better.
We recommend that you go over this page before you make a contribution.

::: tip
You'll need to sign [Uber's CLA](https://cla-assistant.io/uber-go/cff)
before we can accept any of your contributions.
If necessary, a bot will remind
you to accept the CLA when you open your pull request.
:::

## Contributing Ideas

If you'd like to add new functionality to cff,
especially if you want to add any new exported APIs to it,
please [create a new issue](https://github.com/uber-go/cff/issues/new)
describing the problem, your proposed solution,
and anything else that might be relevant.
We will respond to new issues within 30 days with suggestions.
If you're unsure about the direction you want to take with your solution,
feel free to [start a new discussion](https://github.com/uber-go/cff/discussions/new)
instead and we'll assist you if we can.

## Contributing Code

Before contributing code, please follow [these steps](#contributing-ideas).

When your feature is accepted by maintainers, follow these steps:

1. Familiarize yourself with the [Design](#design) of cff.
2. [Test](#testing) your code.
3. Submit a pull request.

We will respond to new PRs with code reviews within 30 days.

## Design

### Phases

cff does its work in two phases: compilation and generation.

- During compilation, the user-provided code is inspected and compiled into an
  internal representation. Any requisite validation is done at this time.
- During generation, the previously built internal representation is used to
  generate code.

### Principles

We tried to follow some standard best practices in designing this tool.

- Do One Thing and Do It Well (from the Unix Philosophy)
- Explicit is better than implicit (from The Zen of Python)

Keeping that in mind, the following guiding principles have helped make
decisions.

- No assumptions should be made, nor restrictions imposed, with regards to code
  organization or abstraction design of user-owned business logic.
  We are only orchestrating functions so we should not need anything except
  function pointers.
- No new dependencies should be added if they can be avoided.
  For example, a dependency on the logging library should not be added if the
  user did not ask for logging.
- No new syntax may be introduced.
  Both, the input and output files must be valid Go code that passes type
  checking.

To support the possibility of online analysis in the future
(versus static analysis)
we want using cff to feel like using a regular Go library.
Users should find it conceivable, based on the APIs and their behavior, that
the functionality of cff is completely runtime.
The presence of code generation should be considered an implementation detail
by users; it should be hidden away until absolutely needed.

To satisfy this,

- The code generator must only make use of information that could be present
  at runtime.
  This means that type information may be used, but not variable names.
  Exceptions may be made here for better UX as long as they don't have
  meaningful behavioral impact.
  For example variable names may be used to affect telemetry, but not graph
  resolution.
- There must be no shared information between Flows unless provided by a
  shared object explicitly.
  For example, Flows cannot share a global task scheduler;
  instead one must be injected at the `cff.Flow` call site.
  Similarly, there must be no shared global configuration between Flows,
  but configuration may be injected at the call site.

## Testing

All logic changes must be accompanied by tests.
We use a healthy mix of unit and integration testing in cff.

For integration testing,
test data is placed inside the internal/tests folder.
Each directory in that folder is treated as its own Go package.
Write sample flows in these directories, generate code from those flows,
and check in the generated code.
Write tests against this generated code.

## Documentation

Before you can contribute documentation, set up your local environment.
Clone the repository, and run:

```bash
cd docs
yarn install
```

Next, run the development server:

```bash
yarn dev
```

Change the Markdown files and preview them in the browser
to make sure everything looks good.

### Formatting

At a high-level, we follow the following styles in Markdown files:

- Use ATX-style headers

  ```markdown
  # Good header

  ---

  Bad header
  ==========
  ```

- Use inline links

  ```markdown
  [Good link](https://example.com)

  ---

  [Bad link][1]

    [1]: https://example.com
  ```

- Use [semantic Line Breaks](https://sembr.org/)

  ```markdown
  Break paragraphs across multiple lines when it makes sense.
  When a sentence ends is always a good place for a line break.
  For long sentences,
  you can do it even in the middle of the sentence,
  to separate out the different clauses.
  ```
