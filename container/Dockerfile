FROM docker.io/library/alpine@sha256:1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870
LABEL maintainer="Ian Spence <ian@ecn.io>"
LABEL org.opencontainers.image.authors="Ian Spence <ian@ecn.io>"
LABEL org.opencontainers.image.source=https://github.com/ecnepsnai/ghpkgsync
LABEL org.opencontainers.image.title="ghpkgsync"
LABEL org.opencontainers.image.description="Github Package Sync"

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk --no-cache update && \
    apk --no-cache add createrepo_c dpkg-dev

ENV GITHUB_USERNAME=
ENV GITHUB_ACCESS_TOKEN=
ENV GITHUB_WEBHOOK_SECRET=
ENV GITHUB_REPOS=
ENV YUM_REPO_ID=
ENV YUM_REPO_DESCRIPTION=
ENV YUM_REPO_BASEURL=

EXPOSE 80 443

VOLUME [ "/repo" ]
WORKDIR "/repo"

ADD ghpkgsync /ghpkgsync
ENTRYPOINT [ "/ghpkgsync" ]