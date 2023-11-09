GO = go
GO_LDFLAGS = -s -w
all: clean merlin
merlin:
	$(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)"
clean:
	rm -f merlin
.PHONY: all clean
