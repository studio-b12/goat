# Builtins

The following builtin functions are available in each script instance.

## `assert`

```ts
function assert(expression: bool, fail_message?: string): void;
```

Takes an `expression` which, when evaluated to `false`, will throw an assert exception. You can pass an additional `fail_message` which will be shown in the exception. This can be used to assert values in responses and fail test execution if they are false.

**Example**

```js
assert(response.StatusCode >= 400, `Status code was ${response.StatusCode}`);
```

## `print`

```ts
function print(...message: string[]): void;
```

Prints the given `message` to the terminal without a leading new line.

**Example**

```js
print("Hello world!");
```

## `println`

```ts
function println(...message: string[]): void;
```

Prints the given `message` to the terminal with a leading new line.

**Example**

```js
println("Hello world!");
```

## `info`

```ts
function info(...message: string[]): void;
```

Logs an *info* log entry to the output logger(s) with the given `message`.

**Example**

```js
info("Hello world!");
```

## `warn`

```ts
function warn(...message: string[]): void;
```

Logs an *warn* log entry to the output logger(s) with the given `message`.

**Example**

```js
warn("Hello world!");
```

## `error`

```ts
function error(...message: string[]): void;
```

Logs an *error* log entry to the output logger(s) with the given `message`.

**Example**

```js
error("Hello world!");
```

## `fatal`

```ts
function fatal(...message: string[]): void;
```

Logs an *fatal* log entry to the output logger(s) with the given `message`. This will also abort the batch execution.

**Example**

```js
fatal("Hello world!");
```

## `debug`

```ts
function debug(...message: string[]): void;
```

Logs an *debug* log entry to the output logger(s) with the given `message`.

**Example**

```js
debug("Hello world!");
```

## `infof`

```ts
function infof(format: string, ...values: any[]): void;
```

Logs an *info* log entry to the output logger(s) with the given `format` formatted with the given `values`. Formatting is handled according to [Go's formatting implementation](https://pkg.go.dev/fmt).

**Example**

```js
infof("Hello %s!", "World");
```

## `warnf`

```ts
function warnf(format: string, ...values: any[]): void;
```

Logs an *warn* log entry to the output logger(s) with the given `format` formatted with the given `values`. Formatting is handled according to [Go's formatting implementation](https://pkg.go.dev/fmt).

**Example**

```js
warnf("Hello %s!", "World");
```

## `errorf`

```ts
function errorf(format: string, ...values: any[]): void;
```

Logs an *error* log entry to the output logger(s) with the given `format` formatted with the given `values`. Formatting is handled according to [Go's formatting implementation](https://pkg.go.dev/fmt).

**Example**

```js
errorf("Hello %s!", "World");
```

## `fatalf`

```ts
function fatalf(format: string, ...values: any[]): void;
```

Logs an *fatal* log entry to the output logger(s) with the given `format` formatted with the given `values`. Formatting is handled according to [Go's formatting implementation](https://pkg.go.dev/fmt). This will also abort the batch execution.

**Example**

```js
fatalf("Hello %s!", "World");
```

## `debugf`

```ts
function debugf(format: string, ...values: any[]): void;
```

Logs an *debug* log entry to the output logger(s) with the given `format` formatted with the given `values`. Formatting is handled according to [Go's formatting implementation](https://pkg.go.dev/fmt).

**Example**

```js
debugf("Hello %s!", "World");
```
