# Project Structure

When you write extensive integration tests for your products, you can use different strategies to re-use often utilized snippets. In this document, you will find some inspiration on how large sets of Goatfiles can be structured.

## Zentralizing commonly used snippets

As you might know, you can extract often used snippets into their own Goatfiles which then can be re-used in other files either via the [`use`](../goatfile/import-statement.md) or via the [`execute`](../goatfile/execute-statement.md) statement.

If you have multiple projects which share similar backend logic, the best approach is to collect shared routines in a separate dependency. For example, you could create a separate Git repository with some utility Goatfiles, which is then added as a Git sub-module to your main projects. Make sure to prefix the name of the target directory in your projects with an underscore (`_`), so that the Goat CLI will not try to execute them by accident. 

Example use-cases for that could be
- logging in with one or more user accounts
- setting default headers for all requests
- encapsulating common procedures like creating and cleaning up entities

## File Structure

In our projects, we employ the following file structure for API integration tests in all our projects. Maybe you can use this for inspiration for your own structure.

```
integrationtests/
├── _shed/
├── issues/
│   └── 123/
│       └── main.goat
├── tests/
│   └── users/
│       ├── _util.goat
│       ├── create.goat
│       ├── list.goat
│       └── delete.goat
├── params.toml.template
├── local.toml
├── staging.toml
└── ci.toml
```

As you can see, we have a directory called `integrationtests/` in all of our projects which contains all tests, test utilities as well as parameter files.

`_shed/` is the name of our dependency containing some utility Goatfiles used i.e. to create new users, log in with users with different permissions, set request defaults and much more. As you can see, this directory is prefixed with `_`, so it will not be executed when calling Goat on the `integrationtests/` directory.

`issues/` contains sub-directories with the name of tickets on our issue board. These are there to demonstrate misbehaviour cases of our API, so these tests should fail on the latest dev state. These can also be used to test against fixes of these issues. When the issue is resolved, these tests should be moved into the respective `tests/` sub-directory.

`tests/` contains the actual integration tests grouped by features. These tests should always pass against the latest dev state, otherwise something might be broken.

We use some user-specific parameters passed into the tests which are stored in different `*.toml` files in the `integrationtests/` directory. These should be specified in the projects `.gitignore`, because every developer might have their own parameters like API keys or user credentials. You could also put a parameter file for automatic tests in there (like the `ci.toml` in our example) which is commited into the repository. The `params.toml.template` is a template file to base custom parameter files on. This is handy because the integration tests expect specific parameters to work with. An example `params.toml.template` could look as following.

```toml
# The server instance to connect to.
instance = "http://localhost:10001"

# Credentials for a user with non-admin privileges.
[credentials.low]
username = "test@test"
password = "password"
apikey = "some api key"

# Credentiasl for a user with admin privileges.
[credentials.admin]
username = "root@root"
password = "password"
apikey = "some api key"
```

## Documentation

To simplify the usage of our Goatfiles, we employ a system of documentation what a Goatfile does, which parameters it expects and – if it is a Goatfile meant to be execuetd or used in other Goatfiles – which state variables it creates which can be further used or captured in an `exec`'s return statement.

The documentation is prefixed with three forward slashes and sits on the very top of the Goatfile. First, stating the purpose of the file followed by required parameters. After that, you can specify the variables created in the execution in the same manner.

Lets take the following as an example for a Goatfile which can be imported in your test to log in two users with two different privilege groups.
```
/// Logs in with a `low` and an `admin` user to the SSP and stores
/// the session in the `default` cookiejar for the `low` user and
/// in the `admin` cookiejar for the `admin` user.
/// 
/// 
/// Required Parameters
/// -----------------------------------------------------------------------
/// instance:                       The root address of the API instance.
/// credentials.low.username:       Username of the `low` user
/// credentials.low.password:       Password of the `low` user
/// credentials.admin.username:     Username of the `admin` user
/// credentials.admin.password:     Password of the `admin` user
/// 
/// 

### Setup

POST {{.instance}}/api/auth/login

[Header]
Content-Type: application/json

[Body]
{
  "username": "{{.credentials.low.username}}",
  "password": "{{.credentials.low.password}}"
}

[Script]
assert(response.StatusCode == 200, `Status code was ${response.StatusCode}`);

---

POST {{.instance}}/api/auth/login

[Options]
cookiejar = "admin"

[Header]
Content-Type: application/json


[Body]
{
  "username": "{{.credentials.admin.username}}",
  "password": "{{.credentials.admin.password}}"
}

[Script]
assert(response.StatusCode == 200, `Status code was ${response.StatusCode}`);

---
```

Here you can find an example for a Goatfile which is created to be used with the `execute` statement in another Goatfile.

```
/// Creates a new user.
/// 
/// 
/// Parameters
/// -----------------------------------------------------------------------
/// instance:   The root address of the API instance.
/// username:   User name of the new user.
/// password:   Password of the new user.
/// 
///
/// Returns
/// -----------------------------------------------------------------------
/// userid:     The ID of the new user.


// ...
```
