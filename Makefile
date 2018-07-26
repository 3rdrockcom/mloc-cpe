# The name of the executable
TARGET=mloc-cpe
.DEFAULT_GOAL: $(TARGET)

# Command arguments and flags
OPTIONS=""

# These are the values we want to pass for VERSION and BUILD
VERSION=1.0.0-dev
BUILD=`git rev-parse HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values
LD_FLAGS=-ldflags="-X github.com/epointpayment/mloc-cpe/app/config.Version=$(VERSION) -X github.com/epointpayment/mloc-cpe/app/config.Build=$(BUILD)"

# Directories required for bindata/embedded data
EMBED=app/migrations/default

# Ignore phony targets
.PHONY: build install clean deps vendor run run-development run-production run-watch

# Builds project
$(TARGET):
	mkdir -p $(EMBED)
	packr build $(LD_FLAGS) -o $(TARGET)

build: $(TARGET)
	@true

# Installs project: copies binary
install:
	packr install $(LD_FLAGS) -o $(TARGET)

# Cleans project: deletes binary
clean:
	if [ -f $(TARGET) ] ; then rm $(TARGET) ; fi

# Get project dependencies
deps:
	go get
	go get github.com/joho/godotenv/cmd/godotenv
	go get github.com/codeskyblue/fswatch
	go get github.com/gobuffalo/packr/...
	go get -u github.com/golang/dep/cmd/dep
	
# Vendor project dependencies
vendor:
	dep ensure

# Runs project: executes binary
run:
	./$(TARGET) ${OPTIONS}

# Runs project: executes binary with development settings
run-development:
	godotenv -f .env.development ./$(TARGET) ${OPTIONS}

# Runs project: executes binary with production settings
run-production:
	godotenv -f .env.production ./$(TARGET) ${OPTIONS}

# Runs project: watches directory for changes and executes binary with development settings
run-watch:
	fswatch -config .fsw.yml