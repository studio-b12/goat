set -e

absolute_path=$(realpath test.goat)
echo "Path to Goatfile: $absolute_path"

goat "$absolute_path"