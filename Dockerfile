FROM golang:1.12 as build
WORKDIR /application
COPY . .
RUN go build

FROM beardedio/terraria
COPY --from=build /application/terraria /tshock/entrypoint
WORKDIR /tshock
ENTRYPOINT [ "/tshock/entrypoint" ]
CMD [ "" ]
