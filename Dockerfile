FROM registry.access.redhat.com/ubi8/ubi-minimal
#FROM registry.access.redhat.com/ubi8/ubi

#ENV TZ="America/New_York" \
ENV LANG="en_US.UTF-8" \
  WEBHOOK=/usr/local/bin/webhook \
  UIDGID=1001:1001

COPY bin/webhook ${WEBHOOK}

USER ${UIDGID}

CMD ["${WEBHOOK}"]
