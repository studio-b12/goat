[VERSION]

# New Features

## Request Defaults [#13]

You can now specify default values for requests which will be applied to each request in the batch when not specified on the request.

*Example:*
```
### Defaults

[Header]
Content-Type: application/json
Accept: application/json

[Script]
assert(response.StatusCode === 200, 
       `Status Code was ${response.StatusCode}`);

### Tests

POST {{.instance}}/objects

[Body]
{
    "name": "Cult of the Lamb",
    "data": {
        "publisher": "Devolver Digital",
        "developer": "Massive Monster",
        "released": "2022-08-11T00:00:00Z",
        "tags": ["Base Building", "Roguelite", "Character Customization"],
        "age_rating": "0"
    }
}

---

// ...
```

## Line of Request in Error Logs [#7]

The file and line where a failed request has been defined is now printed into the log output for a better debugging experience.

*Example:*

> ./test.goat
```
GET https://github.com

[Script]
assert(response.StatusCode === 404);
```

```bash
goat test.goat
```

![](https://github.com/studio-b12/goat/assets/16734205/9abb45e3-c3ad-4702-82df-e84af35c698f)

# Bug Fixes

- Fixed a bug where more than 3 dashes (`-`) as section delimiters after a raw block failed the parsing.
- Removed unnecessary terminal outputs.