#!/bin/bash
set -eu -o pipefail

# parse the current git commit hash
COMMIT=$(git rev-parse HEAD | head -c 7)

# check if the current commit has a matching tag
VERSION=$(git describe --exact-match --abbrev=0 --tags "${COMMIT}" 2> /dev/null || true)

# use the matching tag as the version, if available
if [ -z "$VERSION" ]; then
    VERSION=$COMMIT
fi

# check for changed files (not untracked files)
if [ -n "$(git diff --shortstat 2> /dev/null | tail -n1)" ]; then
    VERSION="${VERSION}-dirty"
fi

echo "$VERSION"
