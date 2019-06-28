#!/usr/bin/env bash

./build.sh toggl-rounder

VERSION=`go run . -version`

printf  "v$VERSION\n" | cat - CHANGELOG.md  > RELEASE_NOTE.md
hub release create --edit -F RELEASE_NOTE.md -a build/toggl-rounder-darwin-amd64 -a build/toggl-rounder-linux-amd64 -a build/toggl-rounder-windows-amd64.exe ${VERSION}
rm RELEASE_NOTE.md
