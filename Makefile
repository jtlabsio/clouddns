BLD_DST=./bin/
BLD_FLGS=-v -a -tags netgo
BNRY_NM=clouddns
CGO_ENABLED=0
DST=${BLD_DST}${BNRY_NM}
GO_CMD=go

build: download
	${GO_CMD} build ${BLD_FLGS} -o ${DST} ./cmd/...

docker:
	docker build --platform linux/amd64 -t clouddns-server .

download:
	${GO_CMD} mod download

.PHONY: build