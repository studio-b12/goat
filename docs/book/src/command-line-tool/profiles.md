# Profiles

Profiles are like parameter files but located in a central file in your home config directory and usable from everywhere.
This should solve the problem when you have similar parameters across multiple projects, so you don't need to duplicate them
everywhere or figure out the path to your centrally stored parameter files.

The location of this file depends on your system type.

| Operating System | Location                                                                                                 |
|------------------|----------------------------------------------------------------------------------------------------------|
| Linux            | `$HOME/.config/goat/profiles*`<br />*or `$XDG_CONFIG_HOME/goat/profiles.*` if `$XDG_CONFIG_HOME` is set* |
| Darwin (OSX)     | `$HOME/Library/Application Support/goat/profiles.*`                                                      |
| Windows          | `%AppData%\goat\profiles.*`                                                                              |

Goat looks for any file in the config directory called `profiles.*`. Supported file types and extensions are as follows.
- `.toml`
- `.yaml` or `.yml`
- `.json`

The structure of a profile file consists of a top level map containing the names of the profiles as keys and the parameter
maps as values.

Here is a quick example.
```yaml
default:
  instance: https://example.com
  credentials:
    token: foobar

staging:
  instance: https://staging.example.com
  credentials:
    token: barbaz

lowprevileges:
  credentials:
    token: lowtoken

local:
  instance: https://localhost:8080
```

Profiles can be addressed when calling `goat` with the `--profile` (or `-P`) parameter. You can define multiple profiles
by chaining multiple of these parameters. For example: 
```
goat -P live -P lowprevileges test.goat
```

The values of the profiles are then merged in order of given parameters.

The `default` profile is **always** loaded into the parameters if present. There you can define parameters that you always
want to use when executing goat. Other passed profiles and parameters will overwrite values in the default profile as usual.

Parameters passed via environment variables or via the `--parameter` flag will also overwrite profile values.
