# Introduction

cff makes it easy to write concurrent code in your application.
Out of the box, it gives you:

- **Bounded goroutines**: cff uses a pool of goroutines to run all operations,
  preventing issues arising from unbounded goroutine growth in your
  application.
- **Panic-safety**: cff prevents panics in your concurrent code from crashing
  your application in a predictable manner.

See [What can I use cff for?](use-cases.md)
to better understand where cff will fit in your application,
or try it with our [Get Started](get-started) section.
