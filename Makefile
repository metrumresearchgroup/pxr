BUILD=`date +%FT%T%z`
tag=`git describe --tags --abbrev=0`
LDFLAGS=-ldflags "-X main.build=${BUILD}"
MAKE_HOME=${PWD}

.PHONY: install experiment cmdtest pxr local-release release

install:
	cd cmd/pxr; go install ${LDFLAGS}

build:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr $(ARGS)

pxr:
	cd cmd/pxr; go build ${LDFLAGS} -o ../../pxr;

local-release:
	cd cmd/pxr; tag=${tag} goreleaser --rm-dist --skip-publish

release:
	cd cmd/pxr; tag=${tag} goreleaser --rm-dist

experiment:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr experiment; rm ./pxr

cmdtest:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr test $(ARGS)

cmdcheck:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr check; rm ./pxr
