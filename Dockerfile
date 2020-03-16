FROM golang:1.14

RUN go get github.com/pseyfert/compile_commands_json_executer

ENTRYPOINT ["/go/bin/compile_commands_json_executer"]
