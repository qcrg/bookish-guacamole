FROM alpine:3.22.1

WORKDIR /var/lib/bookish-guacamole

RUN apk add gcompat

RUN addgroup -S bhgl && adduser bhgl -S bhgl -G bhgl


COPY ./build/bookish-guacamole /bookish-guacamole
COPY ./dev/cert ./dev/cert

RUN chown -R bhgl:bhgl /var/lib/bookish-guacamole

USER bhgl
EXPOSE 8643
CMD ["/bookish-guacamole"]
