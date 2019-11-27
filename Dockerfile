FROM alpine:latest

LABEL maintainer James Rasell<(jamesrasell@gmail.com)> (@jrasell)
LABEL vendor "jrasell"

ENV CHEMTRAIL_VERSION 0.0.1

WORKDIR /usr/bin/

RUN buildDeps=' \
                bash \
                wget \
        ' \
        set -x \
        && apk --no-cache add $buildDeps ca-certificates \
        && wget -O chemtrail https://github.com/jrasell/chemtrail/releases/download/v${CHEMTRAIL_VERSION}/chemtrail_${CHEMTRAIL_VERSION}_linux_amd64 \
        && chmod +x /usr/bin/chemtrail \
        && apk del $buildDeps \
        && echo "Build complete."

ENTRYPOINT ["chemtrail"]

CMD ["--help"]
