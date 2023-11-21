FROM alpine:3.18
WORKDIR /root/
COPY api ./
RUN apk add --no-cache bash tzdata ca-certificates
ENV TZ=Asia/Jakarta
EXPOSE 8080
CMD ["./api"]
