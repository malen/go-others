# Binary name
BINARY=excel
# -w 去掉调试信息
# -s 去掉符号表
#GOBUILD=go build -ldflags '-linkmode "external" -extldflags "-static" -s -w' -o ${BINARY}
GOBUILD=go build -ldflags '-s -w' -o ${BINARY}

GOCLEAN=go clean
RMTARGZ=rm -rf *.gz
VERSION=0.0.1

# Build
.PHONY: build
build:
	$(GOCLEAN)
	$(GOBUILD)

.PHONY: clean
clean:
	$(GOCLEAN)
	$(RMTARGZ)

.PHONY: win
win: clean
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD).exe

.PHONY: release
release:
	# Clean
	$(GOCLEAN)
	$(RMTARGZ)
	# Build for mac
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD)
	tar czvf ${BINARY}-mac64-${VERSION}.tar.gz ./${BINARY}
	# Build for arm
	$(GOCLEAN)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD)
	tar czvf ${BINARY}-arm64-${VERSION}.tar.gz ./${BINARY}
	# Build for linux
	$(GOCLEAN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)
	tar czvf ${BINARY}-linux64-${VERSION}.tar.gz ./${BINARY}
	# Build for win
	$(GOCLEAN)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD).exe
	tar czvf ${BINARY}-win64-${VERSION}.tar.gz ./${BINARY}.exe
	$(GOCLEAN)