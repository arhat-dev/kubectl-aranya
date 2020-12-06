ARG ARCH=amd64

# if we do not use cgo so we can build on alpine and
# copy it to debian
FROM arhatdev/builder-go:alpine as builder
FROM arhatdev/go:debian-${ARCH}
ARG APP=kubectl-aranya-pf

ENTRYPOINT [ "/kubectl-aranya-pf" ]
