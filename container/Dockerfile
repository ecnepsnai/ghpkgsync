FROM docker.io/library/alpine:3
LABEL maintainer="Ian Spence <ian@ecnepsnai.com>"
LABEL org.opencontainers.image.authors="Ian Spence <ian@ecnepsnai.com>"
LABEL org.opencontainers.image.source=https://github.com/ecnepsnai/ghrpmsync
LABEL org.opencontainers.image.title="ghrpmsync"
LABEL org.opencontainers.image.description="Github RPM Sync"

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update && \
    apk add createrepo_c

ENV GITHUB_USERNAME=
ENV GITHUB_ACCESS_TOKEN=
ENV GITHUB_WEBHOOK_SECRET=
ENV GITHUB_REPOS=
ENV YUM_REPO_ID=
ENV YUM_REPO_DESCRIPTION=
ENV YUM_REPO_BASEURL=

EXPOSE 8080

VOLUME [ "/repo" ]
WORKDIR "/repo"

ADD ghrpmsync /ghrpmsync
ENTRYPOINT [ "/ghrpmsync" ]