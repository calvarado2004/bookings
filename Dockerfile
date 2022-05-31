FROM docker.io/bitnami/golang:1.18.1 AS builder
ADD ./ /app
WORKDIR /app
RUN go mod init github.com/calvarado2004/bookings && go mod tidy && go get github.com/alexedwards/scs/v2 && go get github.com/go-chi/chi/v5 && go get github.com/justinas/nosurf && go get github.com/asaskevich/govalidator 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bookings cmd/web/*.go

FROM busybox:stable
ENV APP_HOME /app
RUN adduser 1001 -D -h $APP_HOME && mkdir -p $APP_HOME && chown 1001:1001 $APP_HOME
USER 1001
WORKDIR $APP_HOME
COPY ./templates templates/
COPY ./static static/
COPY --chown=0:0 --from=builder /app/bookings ./
EXPOSE 8080
CMD ["./bookings"]
