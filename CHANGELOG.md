[VERSION]

# ⚠️ Breaking Changes

> Don't worry, these breaking changes only apply to you if you have used Goat in a version before v1.0.0!
> 
> If we introduce breaking changes in following updates, they will be listed here in this way.

- The built-in template function `base64url` has been renamed to `base64Url` for better naming consistency. [#48]

- The built-in template functions `base64` and `base64Url` do now encode strings to Base64 **with leading padding**. If you want to encode without padding, use the new functions `base64Unpadded` and `base64UrlUnpadded`. [#48]

- The previously deprecated alias `headers` for the `header` request section has now been removed.  

- The `BodyJson` field in the `response` object in the `[Script]` section has now been replaced with the `Body` field, which now contains the parsed response body depending on the specified or inferred body type. Read [here](https://studio-b12.github.io/goat/goatfile/requests/script.html) for more information.

- The `BodyRaw` field in the `response` object in the `[Script]` section does now contain the raw response byte array.

# New Features

- Response body parsing does now also support XML responses in addition to JSON. In the Request Option [`responsetype`](https://studio-b12.github.io/goat/goatfile/requests/options.html#responsetype), you can explicitly specify which response type is expected. If not specified, the response body type is inferred via the received `Content-Type` header. The parsed body data can be accessed in the `response.Body` field as addressable object.

- Requesting HTTPS endpoints with self-signed certificates is now possible. [#49, #51]  
  *That means, that the certificate validity is **not checked by default**. You can pass the `--secure` flag or set the `GOATARG_SECURE=true` environment variable to enable certificate validation.*

# Minor Changes and Bug Fixes

- Fixed a bug causing a panic when multiple Goatfiles are passed and one or more of these files do not exist. Now, a proper error message is returned.

- Fixed a bug which caused a panic when invalid key-value paris have been passed via the `--args` flag. [#49]

- When running with the `DEBUG` logging level, the parsed Goatfile is now only printed to log output if the `--dry` flag is passed.

- Some command line parameters now got respective environment variables to set them. These are prefixed with `GOATARG_` and can be viewed in the `--help` overview.