# Defaults Section

## Syntax

> *DefaultsBlock* :  
> `###` `defaults` `NL`+ *PartialRequestProperties*?

## Example

```
### Defaults

[Options]
storecookies = false

[Header]
Authorization: basic {{.token}}

[Script]
assert(response.StatusCode === 200);
```

## Explanation

The `Defaults` section body is structured like a [request](requests/index.md) but without the method and URL section. This includes the blocks `[Options]`, `[Header]`, `[QueryParams]`, `[Body]`, `[PreScript]` and `[Script]`. 

Values specified in `[Header]`, `[Options]` and `[QueryParams]` will be merged with the values set in each request. Values set in a request will overwrite values set in the defaults, if specified in both.

On the other hand, values specified in `[Body]`, `[PreScript]` or `[Script]` will be used if not specified in the requests. If these blocks are specified in the request, they will overwrite the default values as well (even if they are empty).

Default values will also be applied if imported via a `use` directive. Multiple `Defaults` sections will be merged together in order of specification.

### Merge Example

Here is a quick example on how the merges will result.

Let's assume the following defaults are set.
```
### Defaults

[Header]
Content-Type: application/json
Accept-Language: de;q=0.8, en

[Options]
storecookies = false

[Script]
assert(response.Status === 200);
```

Our first request definition looks like this:
```
GET {{.instance}}/api/v1/users/me

[Header]
Authorization: bearer sDrYbXdm2tCYnb8p
```

The resulting **merged** request fields will look as following.
```
GET {{.instance}}/api/v1/users/me

[Header]
Content-Type: application/json
Accept-Language: de;q=0.8, en
Authorization: bearer sDrYbXdm2tCYnb8p

[Options]
storecookies = false

[Script]
assert(response.Status === 200);
```

The second request looks like this.
````
POST {{.instance}}/api/v1/login

[Header]
Content-Type: multipart/form-data; boundary=bbecc22460604a4d97e62895cc6b254a

[Options]
storecookies = true

[Body]
```
--bbecc22460604a4d97e62895cc6b254a
Content-Disposition: form-data; name="username"

foo
--bbecc22460604a4d97e62895cc6b254a
Content-Disposition: form-data; name="password"

bar
--bbecc22460604a4d97e62895cc6b254a--
```

[Script]
assert(response.Status === 200);
assert(response.Body === expectedUserId);
````

And the resulting **merged** parameters will look as following.
````
POST {{.instance}}/api/v1/login

[Header]
Accept-Language: de;q=0.8, en
Content-Type: multipart/form-data; boundary=bbecc22460604a4d97e62895cc6b254a

[Options]
storecookies = true

[Body]
```
--bbecc22460604a4d97e62895cc6b254a
Content-Disposition: form-data; name="username"

foo
--bbecc22460604a4d97e62895cc6b254a
Content-Disposition: form-data; name="password"

bar
--bbecc22460604a4d97e62895cc6b254a--
```

[Script]
assert(response.Status === 200);
assert(response.Body === expectedUserId);
```` 
