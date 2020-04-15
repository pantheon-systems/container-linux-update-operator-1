FROM gcr.io/distroless/base:latest

ARG cmd
COPY bin/${cmd} /bin/${cmd}

ENTRYPOINT ["/bin/${cmd}"]
