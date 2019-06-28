#!/usr/bin/env bash

echo "Building binaries..."
./build.sh toggl-rounder

VERSION=`go run . -version`
echo "Releasing $VERSION"

printf  "v$VERSION\n\n" | cat - CHANGELOG.md  > RELEASE_NOTE.md
hub release create -F RELEASE_NOTE.md -a build/toggl-rounder-darwin-amd64 -a build/toggl-rounder-linux-amd64 -a build/toggl-rounder-windows-amd64.exe ${VERSION}
rm RELEASE_NOTE.md
