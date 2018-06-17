FROM golang:1.10 AS builder

# If we had any external vendor components below would have been a better Dockerfile. This is the real "minimal" one.
# https://medium.com/@pierreprinetti/the-go-dockerfile-d5d43af9ee3c

# Compile the Go code in the Linux container
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM scratch
COPY --from=builder /app ./
COPY ./data/movie_ratings.json ./data/movie_ratings.json
ENTRYPOINT ["./app"]