ARG BUILDER_IMAGE=golang@sha256:4544ae57fc735d7e415603d194d9fb09589b8ad7acd4d66e928eabfb1ed85ff1
ARG RUNTIME_IMAGE=gcr.io/distroless/static@sha256:c9f9b040044cc23e1088772814532d90adadfa1b86dcba17d07cb567db18dc4e


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
