ARG BUILDER_IMAGE=index.docker.io/library/golang@sha256:634cda4edda00e59167e944cdef546e2d62da71ef1809387093a377ae3404df0
ARG RUNTIME_IMAGE=gcr.io/distroless/static@sha256:c9320b754c2fa2cd2dea50993195f104a24f4c7ebe6e0297c6ddb40ce3679e7d


FROM $BUILDER_IMAGE as builder

        WORKDIR /workspace

        COPY .git     .git
        COPY go.mod   go.mod
        COPY go.sum   go.sum
        COPY pkg/     pkg/
        COPY cmd/     cmd/

        RUN set -x && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
                go build -a -v \
			-trimpath \
			-tags osusergo,netgo,static_build \
			-o monero \
				./cmd/monero


FROM $RUNTIME_IMAGE

        WORKDIR /
        COPY --chown=nonroot:nonroot --from=builder /workspace/monero .
        USER nonroot:nonroot

        ENTRYPOINT ["/monero"]
