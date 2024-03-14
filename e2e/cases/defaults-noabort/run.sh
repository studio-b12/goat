res=$(goat main.goat --no-color || true)

if ! echo "$res" | grep 'tests="0/2"'; then
	echo "$res"
	exit 1
fi