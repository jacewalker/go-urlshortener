Inker.ink is a fast web app designed for you to self-host your own URL shortening website on your custom domain.

```
version: "3.8"

services:
  web_app:
    build: .
    container_name: inker.ink
    image: ralphydev/inker.ink:latest
    restart: always
    ports:
      - "8080:8080"
    command: go run main.go
    environment:
      - CUSTOM_DOMAIN="https://inker.ink"
      - SITE_TITLE="inker.ink"
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=short_urls
      - POSTGRES_HOST=database
      - POSTGRES_PORT=5432
    volumes:
      - ./:/go/src
    depends_on:
      - database

  database:
    container_name: go_urlshort_db
    image: postgres:latest
    restart: always
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=short_urls
    ports:
      - 5432:5432
    volumes:
      - db:/var/lib/postgresql/data 

volumes:
  db:

networks:
  default:
      name: inker
```