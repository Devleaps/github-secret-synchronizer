FROM scratch

COPY github-secret-synchronizer /github-secret-synchronizer

ENTRYPOINT ["/github-secret-synchronizer"]
