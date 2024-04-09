go.mod.sri: go.mod go.sum
	OUT=$$(mktemp -d -t nar-hash-XXXXXX) && \
	rm -rf "$$OUT" && \
	go mod vendor -o "$$OUT" && \
	go run tailscale.com/cmd/nardump@v1.62.1 --sri "$$OUT" >go.mod.sri
