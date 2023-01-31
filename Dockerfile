FROM golang:1.19-alpine as builder

WORKDIR /src

COPY . .

RUN go build -o /go/bin/server ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN export APISERVER="https://kubernetes.default.svc"
RUN export SERVICEACCOUNT="/var/run/secrets/kubernetes.io/serviceaccount"
RUN export NAMESPACE="$(cat $SERVICEACCOUNT/namespace)"
RUN export TOKEN="$(cat $SERVICEACCOUNT/token)"
RUN export CACERT="$SERVICEACCOUNT/ca.crt"

COPY --from=builder /go/bin/server /bin/server

EXPOSE 8080

ENTRYPOINT ["/bin/server"]