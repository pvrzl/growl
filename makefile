# The name of the executable (default is current directory name)
TARGET := growl
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 0.0.1
MAINTENERS := ef(***REMOVED***)
# Use linker flags to provide version/build settings to the target
LDFLAGS:=-ldflags "-X main.Version=$(VERSION) -X main.Mainteners=$(MAINTENERS)"

.PHONY: build clean install dep test

dep:
	@go get -u -v github.com/golang/dep/cmd/dep
	@dep ensure -v

build:
	@go build -v $(LDFLAGS) -o $(TARGET) 

clean:
	@rm  $(TARGET)

install:
	@go install $(LDFLAGS)

test:
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out -o cover.html
# @start cover.html
