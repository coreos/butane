FROM registry.fedoraproject.org/fedora:36 AS builder
RUN dnf install -y golang git-core
RUN mkdir /butane
COPY . /butane
WORKDIR /butane
RUN ./build_for_container

FROM registry.fedoraproject.org/fedora-minimal:36
COPY --from=builder /butane/bin/container/butane /usr/local/bin/butane
ENTRYPOINT ["/usr/local/bin/butane"]
