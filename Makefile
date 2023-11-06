CONTRIVED_SERVER_BIN := "bin/contrived"
CONTRIVER_SERVER_SRC := "cmd/contrived/contrived.go"
SRCS := $(shell find . -type f -iname "*.go" -not -path "./.direnv/*")

all: $(CONTRIVED_SERVER_BIN)

$(CONTRIVED_SERVER_BIN): $(SRCS)
	go build -ldflags="-s -w" -o $@ cmd/contrived/contrived.go

.PHONY: clean
clean:
	rm -rf bin/
