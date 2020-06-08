FROM alpine:3.7

RUN apk add --no-cache ca-certificates e2fsprogs findmnt

ADD gluster-simple-csi-plugin /bin/

ENTRYPOINT ["/bin/gluster-simple-csi-plugin"]