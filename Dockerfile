# FROM golang AS download
# WORKDIR /tmp
# # RUN curl -L -o hjm-certcheck https://github.com/hjmcloud/hjm-certcheck/releases/latest/download/hjm-certcheck_linux_$(go env GOARCH)
# ARG src=temp/linux_$(go env GOARCH)
# COPY ${src} ./hjm-certcheck

FROM alpine
# Install dependencies
ARG TARGETARCH
RUN echo "I'm building for $TARGETARCH"
ENV TZ                      Asia/Shanghai
RUN apk update && apk add tzdata ca-certificates bash
# Install hjm-certcheck
ENV WORKDIR                 /app
ADD resource                $WORKDIR/
# ADD ./temp/linux_amd64/main $WORKDIR/main
COPY temp/hjm-certcheck_linux_$TARGETARCH $WORKDIR/hjm-certcheck
RUN chmod +x $WORKDIR/hjm-certcheck

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./hjm-certcheck