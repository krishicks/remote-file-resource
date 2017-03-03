all: assets/in assets/check

clean:
		rm -rf assets

assets:
	  mkdir assets

assets/check: assets cmd/check/main.go
	  CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o assets/check cmd/check/main.go

assets/in: assets cmd/in/main.go
	  CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o assets/in cmd/in/main.go
