FROM golang
RUN apt-get clean && apt-get update
WORKDIR /go/src/github.com/backend_test_RBAC
ADD . .
RUN apt-get install -qy nano
RUN go get -d -v
RUN go install -v
ENTRYPOINT /go/bin/backend_test_RBAC
EXPOSE 8080
