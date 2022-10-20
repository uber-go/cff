# What can't I use cff for?

This list is non-exhaustive,
but we do not recommend trying to use cff for the following use cases:

- **Durable workflow orchestration**:
  cff operates exclusively within the bounds of a single process
  and cannot provide any durability guarantees.
  Instead, use
  [Cadence](https://cadenceworkflow.io/) or [Temporal](https://temporal.io/).
- **Long-running tasks**:
  cff tasks have a definite start and end,
  and cff flows and parallels are typically scoped to a single function call.
  Instead, spawn your own goroutines for that.
  [Carefully](https://github.com/uber-go/guide/blob/master/style.md#dont-fire-and-forget-goroutines).
