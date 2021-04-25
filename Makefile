build-linux:
	env GOOS=linux GOARCH=amd64 go build -o ./bin/coingecko-cache

build-osx:
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/coingecko-cache

install:
	mkdir -p /etc/coingecko-cache
	cp ./bin/coingecko-cache /etc/coingecko-cache/
	cp coingecko-cache.service /etc/systemd/system/
	systemctl daemon-reload 

run:
	./bin/coingecko-cache
