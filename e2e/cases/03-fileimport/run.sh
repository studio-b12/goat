set -e

# CRLF to LF, because git on windows does git on windows things.
# Also '-i' can not be used because sed on OSX does sed on SOX things.
cat body.txt | sed -i 's/^M$//' > body.txt

goat direct.goat
goat imported.goat