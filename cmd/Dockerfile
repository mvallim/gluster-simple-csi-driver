FROM alpine:3.7

LABEL maintainers="Marcos Vallim <tischer@gmail.com>"
LABEL description="GlusterFS Simple Driver"

ENV NODE_ID=""
ENV CSI_ENDPOINT=""

RUN apk add --no-cache ca-certificates e2fsprogs findmnt util-linux

ADD gluster-simple-csi-plugin /bin/

ENTRYPOINT ["/bin/gluster-simple-csi-plugin"]