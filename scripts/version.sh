#!/bin/bash

VERSION_FILE="internal/version/version.go"

VERSION=$(git describe --tags --abbrev=0 2> /dev/null || echo "NA")
COMMIT=$(cut -c-5 <<< "$(git rev-parse HEAD)")
BUILDDATE=$(TZ=UTC date +"%D %H:%M %Z")

cat > $VERSION_FILE <<EOF
package version

const (
	Version    string = "$VERSION"
	CommitHash string = "$COMMIT"
	BuildDate  string = "$BUILDDATE"
)
EOF