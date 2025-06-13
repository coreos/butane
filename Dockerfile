FROM quay.io/fedora/fedora:42 AS builder
RUN dnf install -y golang git-core
RUN mkdir /butane
COPY . /butane
WORKDIR /butane
RUN ./build_for_container

FROM quay.io/fedora/fedora-minimal:42
RUN microdnf install -y tmt python3-pip tmt-provision-container && microdnf clean all
COPY --from=builder /butane/bin/container/butane /usr/local/bin/butane
CMD ["/usr/local/bin/butane"]
