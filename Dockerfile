FROM scratch
RUN mkdir /app
COPY bin/qb_web_server /app/
COPY config /app/
WORKDIR /qb_web_server
CMD ["/qb_web_server"]