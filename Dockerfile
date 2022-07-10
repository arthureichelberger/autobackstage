FROM alpine

LABEL "maintainer"="Arthur Eichelberger <eichelbergerarthur@gmail.com>"
LABEL "repository"="https://github.com/TheDoctor0/zip-release"
LABEL "version"="0.2.0"

COPY autobackstage /opt/bin/autobackstage
RUN chmod +x /opt/bin/autobackstage

RUN ["/opt/bin/autobackstage"]