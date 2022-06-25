# Booking and reservations program in Go

This is the repository for the Bookings and reservations in Go

- Built in Go version 1.18.1
- Uses the [chi browser](github.com/go-chi/chi/v5)
- Uses [alex edwards](github.com/alexedwards/scs/v2) SCS session manager 
- Uses [nosurf](github.com/justinas/nosurf) 

# Admin site

Use the following credentials to check the admin site:

`user: admin@mail.com`

`pass: admin`

![Reservation](./images/reservation.png)

![Calendar](./images/calendar.png)

# Jenkins pipeline for Continuous Integration

This project generates two containers, the init container to make the DB migrations with [Soda](https://gobuffalo.io/documentation/database/migrations/) and the main application container.
\
The Groovy pipeline used to make these containers is available, it uses [Buildah](https://buildah.io/) instead of the whole Docker application

![Jenkins](./images/jenkins.png)

# Docker images generated

The container images are being uploaded to the Docker Hub Public registry

- [init container](https://hub.docker.com/repository/docker/calvarado2004/bookings-init)
- [bookings application](https://hub.docker.com/repository/docker/calvarado2004/bookings)

![Docker Hub](./images/docker-hub.png)



