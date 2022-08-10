FROM golang:1.18 as builder
ADD . /app/
EXPOSE 8080
RUN cd /app && go install github.com/sigmonsays/urlredirector/...
CMD /app/bin/urlredirector
