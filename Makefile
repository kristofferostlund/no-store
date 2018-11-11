build:
	docker build -t no-store-server .

server: build
	docker run -p 4000:4000 --rm \
		no-store-server -port=4000 -address=0.0.0.0

.PHONY: server
