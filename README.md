# kvicksand
Kvicksand is a simple in-memory cache with a basic http rest api interface

All cached values are retained for 30 minutes until they are expired. A scheduled ticker cleans up the expired key-value pairs every 5 seconds.

## How to build

API is containerized and can be run with
```
docker compose up --build
```
Once the API is up and running it will start waiting requests at address
```
0.0.0.0:8080
````
A simple test to see if the API is running from CLI is curling the root
```
curl 0.0.0.0:8080 
```
Which the API should respond by greeting with ```Hello!```

## How to use
API exposes two functional endpoints

* POST 0.0.0.0:8080/<key> with plain text body

Example;
```
curl 0.0.0.0:8080/key1 --data value1
```
* GET 0.0.0.0:8080/<key>

Example;
```
curl 0.0.0.0:8080/key1
```

## Running tests
From the root directory it is possible to run all tests by executing the following command.

```
go test ./... -v -short 
```