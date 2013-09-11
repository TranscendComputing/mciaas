#!/bin/bash
#
# This script only builds the application from source.
set -e

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ];
	do SOURCE="$(readlink "${SOURCE}")"
done
DIR="$(cd -P "$(dirname "${SOURCE}")/.." && pwd)"

# Change into that directory
cd "${DIR}"

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$([ -n "$(git status --porcelain)" ] && echo "+CHANGES" || true )

# If we're building a race-enabled build, then set that up.
if [ ! -z "${MCIAAS_RACE}" ]; then
    echo -e "--> Building with race detection enabled"
    MCIAAS_RACE="-race"
fi

echo -e "--> Installing dependencies..."
go get ./...

# Compile the main MCIaaS app
echo -e "--> Compiling MCIaaS"
go build \
    ${MCIAAS_RACE} \
    -ldflags "-X github.com/TranscendComputing/mciaas/mciaas.GitCommit ${GIT_COMMIT}${GIT_DIRTY}" \
    -v \
    -o mciaas .
