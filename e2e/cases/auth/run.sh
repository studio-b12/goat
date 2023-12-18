set -e

goat direct.goat
goat default.goat
goat params.goat -a "user.name=foo" -a "user.password=bar" -a "token=foobar" -a "type=bearer"