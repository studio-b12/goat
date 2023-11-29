[VERSION]

# ⚠️ Breaking Changes

- The built-in template function `base64url` has been renamed to `base64Url` for better naming consistency. [#48]

- The built-in template functions `base64` and `base64Url` do now encode strings to Base64 **with leading padding**. If you want to encode without padding, use the new functions `base64Unpadded` and `base64UrlUnpadded`. [#48]

# New Features

- Requesting HTTPS endpoints with self-signed certificates is now possible. [#49, #51]  
  *That means, that the certificate validity is **not checked by default**. You can pass the `--secure` flag or set the `GOATARG_SECURE=true` environment variable to enable certificate validation.*

# Minor Changes and Bug Fixes

- Fixed a bug causing a panic when multiple Goatfiles are passed and one or more of these files do not exist. Now, a proper error message is returned.

- Fixed a bug which caused a panic when invalid key-value paris have been passed via the `--args` flag. [#49]

- When running with the `DEBUG` logging level, the parsed Goatfile is now only printed to log output if the `--dry` flag is passed.

- Some command line parameters now got respective environment variables to set them. These are prefixed with `GOATARG_` and can be viewed in the `--help` overview.