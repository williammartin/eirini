ARG stack
FROM packs/cf

ARG buildpacks

WORKDIR /workspace

RUN mkdir -p /packs

COPY builder /packs/
COPY lifecyclebuilder /lifecycle/builder
COPY recipe /packs/

ENTRYPOINT [ \
  "/packs/recipe" \
]

