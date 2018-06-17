# Movie Ratings Stream
This is meant to be a simple golang based service to generate a constant stream of random ratings for random list of movies. 

In order to build this, please use the following command: 

```
  go build -a -installsuffix nocgo -o ./app .
```

When you need to run it:

```
  ./app
```

## REST Endpoint

This is an HTTP streaming endpoint. To connect to the stream:

```
  GET /ratings
```

Returns a constant stream of movie ids with randomized ratings in the below form, separated by newline.

```
  {"data": {"id":680,"rating":9.40509}}
```
