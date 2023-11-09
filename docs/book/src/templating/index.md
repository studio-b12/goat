# Templating

One of the most powerful features of Goat is the ability to use variables from the state in almost every value field in a request definition.

These values are substituted using the powerful templating engine of Go. We will only go over the most basic syntax. Please review the [Go template documentation](https://pkg.go.dev/text/template) for full details.

## Syntax

Template parameters in Goatfiles are defined inside double-braces (`{{ }}`). Variables from the state are referenced using a dot-notation. Let's take a look at the following example.

We have the following state.
```yml
instance: "http://example.com"

credentials:
  username: "foobar"
  password: "password"
```

Now, we can reference the values in our state as following.

```
POST {{.instance}}/api/v1/auth/login

[Body]
{
  "username": "{{.credentials.username}}"
  "password": "{{.credentials.password}}"
}
```

Functions can be called using the name of the function followed by the parameters separated by spaces. Lets take a look at the followign example.

Weassume the following state.
```yaml
name: "Max"
```

```
{{ printf "Hello, %s!" .name }}
```

This will result in the string value `"Hello, Max!"`.

There are also a lot of [builtin functions](./builtins.md) you can use in your Goatfiles.
