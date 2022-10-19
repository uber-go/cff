# CFF

CFF is intended to be a system which facilitates orchestrating large numbers
of interrelated functions with as much concurrent execution as possible.

## Concepts

In CFF, a **Task** is a single executable function or bound method. Tasks have
**inputs** and **outputs**: the parameters and return values of the
corresponding functions. One or more tasks come together to form a **Flow**.
Flows have zero or more **inputs** and one or more **outputs**.

A single Flow must be self-contained. That is, all inputs of all Tasks in a
Flow must come from either another task or as one of the inputs of the Flow
itself.

## Usage

To use CFF, write code using the APIs exported by the
`go.uber.org/cff` package.

```go
var result Response
err := cff.Flow(
    cff.Params(request),
    cff.Results(&response),
    cff.Task(
        client.GetUser),
    cff.Task(
        func(*GetUserResponse) *FooResults {
            ...
        }),
    cff.Task( mapper.FormatResponse),
)
```
