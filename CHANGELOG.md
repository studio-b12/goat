[VERSION]

> **Wran**  
> The option name `[Headers]` is now **deprecated** and support will be removed in release v1.0.0! Please use `[Header]` as option name instead!
>
> Originally, the support for both `Header` and `Headers` as name for the headers option block has been added to avoid errors because both terms have been used somewhat interchangably in the past. But because we want the specification of Goatfiles to be as specific as possible, this support will now be dropped before the official release.

- The definition of the same block multiple times in the same request will now fail the parsing step with an error message alike the following.  
  ```
  2023-03-20T09:35:53Z FATAL execution failed error="failed parsing goatfile at test/test.goat:14:9: [script]: the section has been already defined"
  ```
