FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app
RUN go get "gorm.io/driver/postgres"
RUN go get "gorm.io/gorm"
RUN go get "github.com/gin-gonic/gin"
RUN go build -o main .
EXPOSE 3004
CMD ["/app/main"]


# docker build -t go-urlshort:latest .
# docker run -it -p 4003:8080 -d go-urlshort:latest go run main.go