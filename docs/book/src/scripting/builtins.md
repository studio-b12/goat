# Builtins

The following built-in functions are available in each script instance.

- [`assert`](#assert)
- [`assert_eq`](#assert_eq)
- [`print`](#print)
- [`println`](#println)
- [`info`](#info)
- [`warn`](#warn)
- [`error`](#error)
- [`fatal`](#fatal)
- [`debug`](#debug)
- [`infof`](#infof)
- [`warnf`](#warnf)
- [`errorf`](#errorf)
- [`fatalf`](#fatalf)
- [`debugf`](#debugf)
- [`jq`](#jq)


## `assert`

```ts
function assert(expression: bool, fail_message?: string): void;
```

Takes an `expression` which, when evaluated to `false`, will throw an assert exception. You can pass an additional `fail_message` which will be shown in the exception. This can be used to assert values in responses and fail test execution if they are invalid.

**Example**

```js
assert(response.StatusCode >= 400, `Status code was ${response.StatusCode}`);
```

## `assert_eq`

```ts
function assert_eq(value: any, expected: any, fail_message?: string): void;
```

Takes a `value` and an `expected` value and deep-equals them. That means, when comparing objects and lists, their structure as well as primitive contenst are compared as well. If the comparison fails, it will throw an exception which will display both compared values. You can also pass an additional `fail_message` to further specify the error output.

**Example**

```js
assert_eq(response.StatusCode, 200, "invalid status code");
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

Logs a *warn* log entry to the output logger(s) with the given `message`.

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

Logs a *fatal* log entry to the output logger(s) with the given `message`. This will also abort the batch execution.

**Example**

```js
fatal("Hello world!");
```

## `debug`

```ts
function debug(...message: string[]): void;
```

Logs a *debug* log entry to the output logger(s) with the given `message`.

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

Logs a *warn* log entry to the output logger(s) with the given `format` formatted with the given `values`. Formatting is handled according to [Go's formatting implementation](https://pkg.go.dev/fmt).

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

Logs a *fatal* log entry to the output logger(s) with the given `format` formatted with the given `values`. Formatting is handled according to [Go's formatting implementation](https://pkg.go.dev/fmt). This will also abort the batch execution.

**Example**

```js
fatalf("Hello %s!", "World");
```

## `debugf`

```ts
function debugf(format: string, ...values: any[]): void;
```

Logs a *debug* log entry to the output logger(s) with the given `format` formatted with the given `values`. Formatting is handled according to [Go's formatting implementation](https://pkg.go.dev/fmt).

**Example**

```js
debugf("Hello %s!", "World");
```

## `jq`

```ts
function jq(object: any, src: string): any[];
```

For more complex analysis on JSON responses, `jq` can be used to perform JQ commands on any passed `object`. Goat uses [itchyny/gojq](https://github.com/itchyny/gojq) as implementation of JQ. [Here](https://jqlang.github.io/jq/manual) you can find more information about the JQ syntax.

The results are always returned as a list. If there are no results, an empty list is returned. When the command compilation fails or the JQ execution fails, the function will throw an exception.

**Example**

*Explanation: Lists all instances recursively where the value of `href` is an object and is empty.*

```js
jq(response.Body, `..
    | if type == "object" then .href else empty end
    | if type == "object" then . else empty end
    | select( . | length == 0 )`);
```
