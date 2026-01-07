## Build Cross platform installer

GOOS=windows GOARCH=amd64 go build -o build/printer-agent.exe
GOOS=darwin  GOARCH=amd64 go build -o build/printer-agent
GOOS=linux   GOARCH=amd64 go build -o build/printer-agent


## Build and installer for windows

docker run --rm -i -v "$PWD:/work" amake/innosetup installer.iss