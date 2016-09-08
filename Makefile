all:
	cgogen pm.yml

clean:
	rm -f pm/cgo_helpers.go pm/cgo_helpers.h pm/cgo_helpers.c
	rm -f pm/doc.go pm/types.go pm/const.go
	rm -f pm/pm.go

test:
	cd pm && go build
	