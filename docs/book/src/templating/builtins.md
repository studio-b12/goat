# Builtins

The following builtin functions are available in templates used in Goatfiles.

## `base64`

```
base64 <value: string> -> string
```

Returns the input `value` as base64 encoded string.

**Example:**

```
{{ base64 "hello world" }}
```

## `base64url`

```
base64url <value: string> -> string
```

Returns the input `value` as base64 URL encoded string.

**Example:**

```
{{ base64url "hello world" }}
```

## `md5`

```
md5 <value: string> -> string
```

Returns the HEX encoded MD5 hash of the given input `value`.

**Example:**

```
{{ md5 "hello world" }}
```

## `sha1`

```
sha1 <value: string> -> string
```

Returns the HEX encoded SHA1 hash of the given input `value`.

**Example:**

```
{{ sha1 "hello world" }}
```

## `sha256`

```
sha256 <value: string> -> string
```

Returns the HEX encoded SHA256 hash of the given input `value`.

**Example:**

```
{{ sha256 "hello world" }}
```

## `sha512`

```
sha512 <value: string> -> string
```

Returns the HEX encoded SHA512 hash of the given input `value`.

**Example:**

```
{{ sha512 "hello world" }}
```

## `randomString`

```
randomString <length?: integer> -> string
```

Returns a random string with the given length. If no length is passed, the default length is 8 characters.

**Example:**

```
{{ randomString 16 }}
```

## `randomInt`

```
randomInt <n?: integer> -> int
```

Returns a random integer in the range `[0, n)` where `n` is given as parameter. if no parameter is passed, `n` defaults to the max int value.

**Example:**

```
{{ randomInt 256 }}
```

## `timestamp`

```
timestamp <format?: string> -> string
```

Returns the current timestamo in the given format. The format is specified in [Go's time format](https://pkg.go.dev/time). If no format is given, the time is returned as Unix seconds string.

**Example:**

```
{{ timestamp "Mon, 02 Jan 2006 15:04:05 MST" }}
```

## `isset`

```
timestamp <map: map[string]any> <key: string> -> bool
```

Returns `true` when the given `key` is present and its corresponding value is not `nil` in the given `map`. Otherwise, `false` is returned.

**Example:**

```
{{ isset . "username" }}
```
