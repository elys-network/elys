# Coding Guidelines

This document is an extension to [CONTRIBUTING](./CONTRIBUTING.md) and provides more details about the coding guidelines and requirements.

## API & Design

- Code must be well structured:
  - packages must have a limited responsibility (different concerns can go to different packages),
  - types must be easy to compose,
  - think about maintainbility and testability.
- "Depend upon abstractions, [not] concretions".
- Try to limit the number of methods you are exposing. It's easier to expose something later than to hide it.
- Take advantage of `internal` package concept.
- Follow agreed-upon design patterns and [naming conventions](https://medium.com/@kdnotes/golang-naming-rules-and-conventions-8efeecd23b68).
- publicly-exposed functions are named logically, have forward-thinking arguments and return types.
- Avoid global variables and global configurators.
- Favor composable and extensible designs.
- Minimize code duplication.
- Limit third-party dependencies.

Performance:

- Avoid unnecessary operations or memory allocations.

Security:

- Pay proper attention to exploits involving:
  - gas usage (avoid continuously growing gas fees per operation)
  - transaction verification and signatures
  - malleability
  - code must be always deterministic (possible non-deterministic operations like random, sort, functions that behave differently per OS)
- Thread safety. If some functionality is not thread-safe, or uses something that is not thread-safe, then clearly indicate the risk on each level. (Some operations could have different result if it's executed in parallel e.g. swap operations in parallel)

## Automated Tests

Make sure your code is well tested:

- Provide unit tests for every unit of your code if possible. Unit tests are expected to comprise 70%-80% of your tests.
- Describe the test scenarios you are implementing for integration tests.
- Create integration tests for queries and msgs.
- Use both test cases and property / fuzzy testing. We use the [rapid](pgregory.net/rapid) Go library for property-based and fuzzy testing.
- Do not decrease code test coverage. Explain in a PR if test coverage is decreased.

We expect tests to use `require` or `assert` rather than `t.Skip` or `t.Fail`,
unless there is a reason to do otherwise.
When testing a function under a variety of different inputs, we prefer to use
[table driven tests](https://github.com/golang/go/wiki/TableDrivenTests).
Table driven test error messages should follow the following format
`<desc>, tc #<index>, i #<index>`.
`<desc>` is an optional short description of whats failing, `tc` is the
index within the test case table that is failing, and `i` is when there
is a loop, exactly which iteration of the loop failed.
The idea is you should be able to see the
error message and figure out exactly what failed.
Here is an example check:

```go
<some table>
for tcIndex, tc := range cases {
  <some code>
  resp, err := doSomething()
  require.NoError(err)
  require.Equal(t, tc.expected, resp, "should correctly perform X")
```
