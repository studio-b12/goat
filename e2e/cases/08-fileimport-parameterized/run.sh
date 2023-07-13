set -e

# CRLF to LF, because git on windows does git on windows things.
# Also '-i' can not be used because sed on OSX does sed on SOX things.
cat body.txt | sed 's/^M$//' | tee body.txt

goat -a "file=body.txt" direct.goat
goat -a "file=body.txt" imported.goat