---
name: Bug report
about: Report a problem with the SDK
labels: bug
---

**SDK version**
e.g. v0.1.0, or the commit SHA if not using a tagged release.

**Go version**
Output of `go version`.

**Describe the bug**
A clear description of what went wrong, including the exact error message
or `APIError` fields if applicable.

**Minimal reproduction**

```go
// A minimal, self-contained code sample that reproduces the issue.
```

**Expected behavior**
What you expected to happen instead.

**Additional context**
Anything else relevant: retry policy in use, whether the request was
concurrent, whether a custom http.Client or observer was configured, etc.
