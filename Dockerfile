FROM golang:latest AS builder
RUN mkdir /fcct
COPY . /fcct
WORKDIR /fcct
RUN ./build_releases

FROM scratch
COPY --from=builder /fcct/bin/releases/fcct-x86_64-unknown-linux-gnu /usr/local/bin/fcct
ENTRYPOINT ["/usr/local/bin/fcct"]
