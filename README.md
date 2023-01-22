# 🧞‍♀️ gurl

A CLI tool to simplify and automate integration testing of HTTP APIs by using script files.

> Example Gurlfile
```
use util/login

### Setup

LOGIN {{.instance}}/api/v1/auth

[Header]
Content-Type: application/json

[Body]
{ 
    "username": "{{.username}}",
    "password": "{{.password}}"
}

[Script]
assert(response.StatusCode == 200, `Status code was ${response.StatusCode}`);

---

### Tests

GET {{.instance}}/api/v1/list

[Script]
assert(response.StatusCode == 200, `Status code was ${response.StatusCode}`);
print(response.Body);
```

## Getting Started

### Installation

If you have the Go toolchain installed, you can simply install the tool via `go install`.
```
go install github.com/studio-b12/gurl/cmd/gurl
```

Otherwise, you can also download the binaries from the [Releases Page](https://github.com/studio-b12/gurl/releases).

### Gurlfile

Now, you can dive in to create your first Gurlfile. You can find more information on how Gurlfiles work here:

- [Gurlfile Specification](docs/gurlfile-spec.md)
- [Implementation Docs](docs/implementation.md)

You can also simply generate an example Gurlfile using the following command.
```
gurl --new
```

There, you can define your setup, teardown and test requests as you desire.

## Contribute

If you find any issues, want to submit a suggestion for a new feature or improvement of an existing one or just want to ask a question, feel free to [create an Issue](https://github.com/studio-b12/gurl/issues/new).

If you want to contribute to the project, just [create a fork](https://github.com/studio-b12/gurl/fork) and [create a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request) with your changes. We are happy to review your contribution and make you a part of the project. 😄

---

© 2023 Studio B12 GmbH / B12-Touch GmbH  
https://studio-b12.de / https://b12-touch.de

Covered by the MIT License.