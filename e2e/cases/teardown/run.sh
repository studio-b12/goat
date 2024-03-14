res=$(goat test.goat --silent || true)

if [ "$res" != "123" ]; then
	echo "$res != 123"
	exit 1
fi