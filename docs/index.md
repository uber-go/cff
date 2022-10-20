---
tagline: Concurrency toolkit for Go
home: true
footer: Made by Uber with ❤️
features:
    - details: |
        cff helps you avoid common bugs in concurrent code by isolating and
        hiding complexity behind easy to use APIs.
      title: Easy
    - details: |
        cff moves expensive operations to compile time so that at runtime,
        the only cost your code pays is to run the functions themselves.
      title: Fast
    - details: |
        cff validates your code for type-safety at compile time,
        and generates code to safely execute those operations at runtime.
      title: Safe
actionText: Get Started →
actionLink: /get-started/
---

cff (pronounce *caff* as in caffeine) is a library and code generator for Go
that makes it easy to run independent and interdependent Go functions
concurrently.

Check out [What can I use cff for?](use-cases.md) to understand places where
cff might fit in your application.
