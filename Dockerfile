FROM quay.io/fedora/fedora:42 AS builder
RUN dnf install -y golang git-core
RUN mkdir /butane
COPY . /butane
WORKDIR /butane
RUN ./build_for_container

FROM quay.io/fedora/fedora-minimal:42
COPY --from=builder /butane/bin/container/butane /usr/local/bin/butane
ENTRYPOINT ["/usr/local/bin/butane"]
