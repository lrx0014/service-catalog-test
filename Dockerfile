FROM scratch

COPY ./pkg/catalog/main /main
ENTRYPOINT [ "/main" ]