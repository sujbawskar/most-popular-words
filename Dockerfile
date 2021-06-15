FROM golang:1.16 as builder

RUN mkdir -p /src/cmd
WORKDIR /src/cmd
COPY . /src/cmd
RUN make build && cp platform-go-tech-test /usr/bin/

FROM alpine:3.9
COPY --from=builder /usr/bin/platform-go-tech-test /usr/bin/
COPY --from=builder /src/cmd/world192.txt /usr/bin/
WORKDIR /usr/bin/
ENTRYPOINT ["platform-go-tech-test"]