.PHONY: all
all: cyr2lat cyr2lat_win.exe

cyr2lat: main.go
	GOOS=linux GOARCH=amd64 go build -o cyr2lat main.go

cyr2lat_win.exe: main.go
	GOOS=windows GOARCH=amd64 go build -o cyr2lat_win.exe main.go

clean:
	rm -rf cyr2lat cyr2lat_win.exe
