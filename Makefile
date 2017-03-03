all: assets/in assets/check

clean:
		rm -rf assets

assets:
	  mkdir assets

assets/check: assets cmd/check/main.go
	  GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o assets/check cmd/check/main.go

assets/in: assets cmd/in/main.go
	  GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o assets/in cmd/in/main.go
