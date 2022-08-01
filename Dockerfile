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