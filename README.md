## Build Cross platform installer

GOOS=windows GOARCH=amd64 go build -o build/printer-agent.exe
GOOS=darwin  GOARCH=amd64 go build -o build/printer-agent
GOOS=linux   GOARCH=amd64 go build -o build/printer-agent

## Build to run in background

GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o printer-agent.exe


## Build and installer for windows

docker run --rm -i -v "$PWD:/work" amake/innosetup installer.iss

## Create certificate

mkcert 127.0.0.1 localhost

For build you also need rootCA.pem file which can be found here

``` bash
mkcert -CAROOT
```

In mac there might be issue printing on process restart, Run the below command to reset the connection:

``` bash
sudo pkill bluetoothd
```