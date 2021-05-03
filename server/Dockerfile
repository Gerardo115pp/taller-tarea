FROM golang:alpine3.13

WORKDIR /home/el_maligno/Documents/OctavoSemestre/SemEstructurasII/ActividadV/server
COPY . .


RUN go get -d -v ./...
RUN go install -v ./...

CMD ["taller-server"]