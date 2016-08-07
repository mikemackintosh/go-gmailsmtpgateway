all: test build

setup:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update

lint:
	echo "gometalinter with vet, golint, gofmt, deadcode, and ineffassign..."
	gometalinter --deadline=60s \
	             --disable-all --enable=vet --enable=golint --enable=gofmt --enable=deadcode --enable=ineffassign \
	             --exclude=vendor ./...

vet:
	go list ./... | grep -v vendor | xargs go vet

test: lint vet
	go test ./...
	
build:
	go build -o bin/gmailsmtpd cmd/gmailsmtpd.go
