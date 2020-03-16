FROM golang:1.14

RUN go get github.com/pseyfert/extractor_action/extractor_action_cmd

ENTRYPOINT ["/go/bin/extractor_action_cmd"]
