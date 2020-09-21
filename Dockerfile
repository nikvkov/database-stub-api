FROM golang:latest
WORKDIR /go/src/github.com/nikvkov/database-stub-api
COPY . .
#RUN go mod init
RUN go build -o stub-server database_stub_api.go
RUN ls
CMD ./stub-server