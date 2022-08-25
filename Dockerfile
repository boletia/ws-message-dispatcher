# syntax=docker/dockerfile:1
FROM golang:1.14-alpine 

RUN mkdir -p {/app,/app/build,/app/scripts}
WORKDIR /app
ADD . ./

RUN apk update
RUN apk add git make bash
RUN apk --no-cache add ca-certificates git


# RUN git config --global --add url."https://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@github.com/boletia/".insteadOf "https://github.com/boletia/" \
#     && GOPRIVATE=$GOPRIVATE \
#     && export GOOS=$linux \
#     && export GOARCH=amd64

#RUN git config --global url."https://541e96dbb1b3edc724efe63828d0c2a01d75a5c0:x-oauth-basic@github.com/".insteadOf "https://github.com/"


CMD ["/bin/sh"]