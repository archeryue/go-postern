local:
	cd cmd/local/ && go build -o ../../build/local local.go

remote:
	cd cmd/remote && go build -o ../../build/remote remote.go

.PHONY : clean

clean:
	rm -f build/local
	rm -f build/remote
