FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-oracle-fusion-cloud"]
COPY baton-oracle-fusion-cloud /