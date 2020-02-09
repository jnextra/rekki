FROM golang:alpine AS builder
RUN apk update && apk add --no-cache make git libcap && rm -rf /var/cache/apk/*
RUN adduser -D -g '' appuser
WORKDIR $GOPATH/src/github.com/josephn123/rekki
COPY Makefile Makefile
RUN make bootstrap
COPY cmd cmd
COPY pkg pkg
RUN CGO_ENABLED=0 make && \
    setcap 'cap_net_bind_service=+ep' dist/rekki && \
    chown -R appuser dist/rekki && \
    mv dist/rekki /go/bin/rekki

FROM alpine
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/rekki /rekki/rekki
USER appuser
EXPOSE 8080
ENTRYPOINT ["/rekki/rekki"]