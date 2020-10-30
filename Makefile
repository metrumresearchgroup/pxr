BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.build=${BUILD}"
MAKE_HOME=${PWD}

.PHONY: install experiment cmdtest pxr

install:
	cd cmd/pxr; go install ${LDFLAGS}

build:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr $(ARGS)

pxr:
	cd cmd/pxr; go build ${LDFLAGS} -o ../../pxr;

experiment:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr experiment; rm ./pxr

cmdtest:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr test $(ARGS)

cmdcheck:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr check; rm ./pxr
