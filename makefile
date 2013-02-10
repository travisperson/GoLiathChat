all: client server

CLIFILES=ccgClient.go ccgPacket.go ccgHost.go ccgTermbox.go ccgTools.go ccgUser.go
client: $(CLIFILES)
	go build -o client $(CLIFILES)

SERVFILES=ccgServer.go ccgPacket.go ccgTools.go ccgUser.go
server: $(SERVFILES)
	go build -o server $(SERVFILES)
