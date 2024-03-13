# Built-ins

The following built-in functions are available in templates used in Goatfiles.

- [`base64`](#base64)
- [`base64Url`](#base64Url)
- [`base64Unpadded`](#base64Unpadded)
- [`base64UrlUnpadded`](#base64UrlUnpadded)
- [`formatTimestamp`](#formatTimestamp)
- [`md5`](#md5)
- [`sha1`](#sha1)
- [`sha256`](#sha256)
- [`sha512`](#sha512)
- [`randomString`](#randomString)
- [`randomInt`](#randomInt)
- [`timestamp`](#timestamp)
- [`isset`](#isset)
- [`json`](#json)

## `base64`

```
base64 <value: string> -> string
```

Returns the input `value` as base64 encoded string with padding.

**Example:**

```
{{ base64 "hello world" }}
```

## `base64Url`

```
base64Url <value: string> -> string
```

Returns the input `value` as base64url encoded string with padding.

**Example:**

```
{{ base64Url "hello world" }}
```

## `base64Unpadded`

```
base64Unpadded <value: string> -> string
```

Returns the input `value` as base64 encoded string without padding.

**Example:**

```
{{ base64Unpadded "hello world" }}
```

## `base64UrlUnpadded`

```
base64UrlUnpadded <value: string> -> string
```

Returns the input `value` as base64url encoded string without padding.

**Example:**

```
{{ base64UrlUnpadded "hello world" }}
```

## `formatTimestamp`

```
formatTimestamp <value: Date> <format?: string> -> string
formatTimestamp <value: string> <valueFormat: string> <format?: string> -> string
```

Takes either a `Date` object and an optional output format or a timestamp string, the format for the input timestamp and
an optional output format and returns the formatted timestamp.

The format is according to Go's [`time` package definition](https://pkg.go.dev/time#pkg-constants). You can also specify
the names of the predefined formats like `rfc3339` or `DateOnly`. If no format is passed, the time will be represented as
Unix seconds.

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

Returns the current timestamp in the given format. 

The format is according to Go's [`time` package definition](https://pkg.go.dev/time#pkg-constants). You can also specify
the names of the predefined formats like `rfc3339` or `DateOnly`. If no format is passed, the time will be represented as
Unix seconds.

**Example:**

```
{{ timestamp "Mon, 02 Jan 2006 15:04:05 MST" }}
```

## `isset`

```
isset <map: map[string]any> <key: string> -> bool
```

Returns `true` when the given `key` is present and its corresponding value is not `nil` in the given `map`. Otherwise, `false` is returned.

**Example:**

```
{{ isset . "username" }}
```

## `json`

```
json <value: any> <ident?: string | int> -> string
```

Serializes a given `value` into a JSON string. You can pass a string value used as ident or a number of spaces used as indent.

**Example:**

```
{{ json .someObject 2 }}
```
