# This how we want to name the binary output
BINARY=`basename \`pwd\``

# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"
VERSION=`git describe --tags|sed -e "s/\-/\./g"`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
TAGS=-tags ""
CMD_NAMESPACE=$(shell go list ./cmd)
BRANCH := $(shell git show-ref | grep "$(REVISION)" | grep -v HEAD | awk '{print $$2}' | sed 's|refs/remotes/origin/||' | sed 's|refs/heads/||' | sort | head -n 1)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS_DEB=-ldflags=all="-X ${CMD_NAMESPACE}.Version=${VERSION}.Version -X ${CMD_NAMESPACE}.Date=${BUILD} -X ${CMD_NAMESPACE}.Commit=${COMMIT} -X ${CMD_NAMESPACE}.Branch=${BRANCH}"
LDFLAGS_REL=-ldflags=all="-w -s -X ${CMD_NAMESPACE}.Version=${VERSION} -X ${CMD_NAMESPACE}.Date=${BUILD} -X ${CMD_NAMESPACE}.Commit=${COMMIT} -X ${CMD_NAMESPACE}.Branch=${BRANCH}"


# Builds the project
lrelease:
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build ${TAGS} ${LDFLAGS_REL} -o ${BINARY}
release:
	env CGO_ENABLED=0 go build ${TAGS} ${LDFLAGS_REL} -o ${BINARY}
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build ${TAGS} ${LDFLAGS_REL} -o ${BINARY}

ldebug:
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build ${TAGS} ${LDFLAGS_DEB} -o ${BINARY}

debug:
	env CGO_ENABLED=0 go build ${TAGS} ${LDFLAGS_DEB} -o ${BINARY}
	env GOOS=linux GOARCH=amd64 go build ${TAGS} ${LDFLAGS_DEB} -o ${BINARY}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
