# Getting Started

## Installation

### Pre-Built Binaries 

First of all, you need to install the `goat` tool. It is distributed as a single binary without any dependencies. So simply download the [**latest release binaries from GitHub**](https://github.com/studio-b12/goat/releases) fitting your system configuration.

You can also use the provided [installation script](https://github.com/studio-b12/goat/blob/main/scripts/download.sh), if you want.

```bash
$ curl -Ls https://raw.githubusercontent.com/studio-b12/goat/main/scripts/download.sh | bash -
```

### Installation with `go install`

Alternatively, if you have the Go toolchain installed on your system, you can also simply use `go install` to build and install the tool for your system. This is also very useful if no pre-built binary is available for your specific system configuration.

> Details on how to install the Go toolchain is available [here](https://go.dev/doc/install).

Now, simply use the following command to install the program.
```
go install github.com/studio-b12/goat/cmd/goat@latest
```

This will install the latest tagged version. You can also specify a specific version or branch after the `@`.
 ```
go install github.com/studio-b12/goat/cmd/goat@dev
```

## First Goatfile

After that, simply use the following command to initalize a new Goatfile.
```
goat --new
```

This will generate a new Goatfile with the name `tests.goat` at your current working directory with some simple examples and documentation. You can use this file as your starting point to play around with Goat or to create your projects integration test structure.
