FROM busybox
MAINTAINER gyr666
WORKDIR /
COPY ./concurrentNet /main
EXPOSE 9090
ENTRYPOINT ["./main"]
