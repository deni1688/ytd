FROM alpine

WORKDIR /app

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY ./server.crt /usr/local/share/ca-certificates/server.crt

RUN update-ca-certificates
RUN apk add --no-cache ffmpeg
RUN mkdir public
RUN mkdir downloads

COPY public ./public 
COPY downloads ./downloads 
COPY ytd .

CMD ./ytd

