tidy:
	go mod tidy; go mod vendor;

run:
	. ./dev.rc && (cd cmd; go build . && ./cmd)

prod:
	. ./prod.rc && (cd cmd; go build . && ./cmd)