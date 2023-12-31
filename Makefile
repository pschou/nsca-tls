VERSION = 0.1.$(shell date +%Y%m%d.%H%M)
FLAGS := "-s -w -X main.version=${VERSION}"


build:
	CGO_ENABLED=0 go build -ldflags=${FLAGS} -o nsca-tls-client client.go lib*.go
	CGO_ENABLED=0 go build -ldflags=${FLAGS} -o nsca-tls-server server.go lib*.go
	CGO_ENABLED=0 go build -ldflags=${FLAGS} -o nsca-tls-post post.go lib*.go
	upx --lzma nsca-tls-client
	upx --lzma nsca-tls-server
	upx --lzma nsca-tls-post
