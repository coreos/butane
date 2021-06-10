FROM registry.fedoraproject.org/fedora:34 AS builder
RUN dnf install -y golang git
RUN mkdir /butane
COPY . /butane
WORKDIR /butane
RUN ./build_for_container

FROM scratch
COPY --from=builder /butane/bin/container/butane-x86_64-unknown-linux-gnu /usr/local/bin/butane
ENTRYPOINT ["/usr/local/bin/butane"]
