BIN = ssh-login-notify

$(BIN):
	go build -o $(BIN) --ldflags "-s -w" .

clean:
	rm -f $(BIN)