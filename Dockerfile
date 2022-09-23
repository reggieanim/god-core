FROM ghcr.io/go-rod/rod

WORKDIR /build


COPY cleanNotion .
COPY cleanNotion.json .

CMD ["./cleanNotion"]