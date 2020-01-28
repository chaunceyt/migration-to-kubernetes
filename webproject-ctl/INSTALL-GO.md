# Installing go


```
GOOS=linux
GO_VERSION: 1.13.1
curl -OL https://storage.googleapis.com/golang/go${GO_VERSION}.${GOOS}-amd64.tar.gz
tar -xf go${GO_VERSION}.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo mv go /usr/local
mkdir -p "$HOME/go/bin"
rm -rf go${GO_VERSION}.linux-amd64.tar.gz
go version
```
