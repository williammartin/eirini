ARG stack
FROM golang as lifecycle
ARG diego_version=2.30.0
ARG diego_repo=github.com/cloudfoundry/diego-release
ARG bal_repo=code.cloudfoundry.org/buildpackapplifecycle

WORKDIR /diego
RUN git clone --branch "v${diego_version}" --depth 1 "https://${diego_repo}" . && \
  git submodule update --init --recursive
RUN GOPATH=/diego CGO_ENABLED=0 go install -a -installsuffix static "${bal_repo}/..."

FROM packs/cf

ARG buildpacks

WORKDIR /workspace

RUN mkdir -p /packs

COPY --from=lifecycle /diego/bin/builder /packs/
COPY --from=lifecycle /diego/bin /lifecycle
COPY recipe /packs/

ENTRYPOINT [ \
  "/packs/recipe" \
]

