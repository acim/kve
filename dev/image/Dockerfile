FROM golang:1.16.0-buster

RUN curl -sL https://deb.nodesource.com/setup_14.x | bash - && \
    apt-get install -y nodejs

RUN groupadd -r -g 1000 acim \
    && useradd -r -u 1000 -g 1000 -m -c "developer" -s /bin/bash acim

RUN mkdir /app && \
    chown acim:acim /app

USER 1000

RUN go get -u github.com/cespare/reflex

COPY --chown=acim:acim reflex.conf /home/acim/etc/
COPY --chown=acim:acim build.sh start.sh /home/acim/bin/

WORKDIR /app

ENTRYPOINT ["/home/acim/bin/start.sh"]

EXPOSE 3000 35729
