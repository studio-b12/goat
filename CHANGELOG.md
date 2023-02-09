[VERSION]

- Fixed runtime panic when executing a request which has no `[Script]` section.
- Fixed a bug where headers could not be analyzed in scripts.
- Added template functions `randomString` and `randomInt`. [Here](https://github.com/studio-b12/goat/blob/main/docs/implementation.md#parameters) you can find more details. [#15]
- Added template function `timestamp`. [Here](https://github.com/studio-b12/goat/blob/main/docs/implementation.md#parameters) you can find more details. [#16]