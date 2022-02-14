FROM registry.fedoraproject.org/fedora:35 AS builder
ARG BUILDARCH
ENV BUILDARCH=${BUILDARCH}
RUN dnf install -y golang git
RUN mkdir /butane
COPY . /butane
WORKDIR /butane
RUN ./build_for_container

FROM scratch
ARG BUILDARCH
ENV BUILDARCH=${BUILDARCH}
COPY --from=builder /butane/bin/container/butane-${BUILDARCH}-unknown-linux-gnu /usr/local/bin/butan
ENTRYPOINT ["/usr/local/bin/butane"]
