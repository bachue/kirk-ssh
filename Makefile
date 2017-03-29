.PHONY: ALL clean

ALL=kirk-ssh.darwin.amd64 kirk-ssh.linux.386 kirk-ssh.linux.amd64
SOURCES=main.go glide.lock

ALL: $(ALL)

kirk-ssh.darwin.amd64: $(SOURCES)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o kirk-ssh.darwin.amd64 .
kirk-ssh.linux.386: $(SOURCES)
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o kirk-ssh.linux.386 .
kirk-ssh.linux.amd64: $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o kirk-ssh.linux.amd64 .
clean:
	rm -f $(ALL)
