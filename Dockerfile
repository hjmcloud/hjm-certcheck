FROM golang AS download
WORKDIR /tmp
RUN curl -L -o hjm-certcheck https://github.com/hjmcloud/hjm-certcheck/releases/latest/download/hjm-certcheck_linux_$(go env GOARCH)

FROM alpine
# Install dependencies
ENV TZ                      Asia/Shanghai
RUN apk update && apk add tzdata ca-certificates bash
# Install hjm-certcheck
ENV WORKDIR                 /app
ADD resource                $WORKDIR/
# ADD ./temp/linux_amd64/main $WORKDIR/main
COPY --from=download /tmp/hjm-certcheck $WORKDIR/hjm-certcheck
RUN chmod +x $WORKDIR/hjm-certcheck

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./hjm-certcheck