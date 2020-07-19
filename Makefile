BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.build=${BUILD}"
MAKE_HOME=${PWD}

.PHONY: install experiment cmdtest pxr

install:
	cd cmd/pxr; go install ${LDFLAGS}

pxr:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr;

experiment:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr experiment; rm ./pxr

cmdtest:
	cd cmd/pxr; go build ${LDFLAGS} -o pxr; ./pxr test; rm ./pxr
