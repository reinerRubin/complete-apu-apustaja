# completer helper

## basic build and launch
```
$ make build-docker
...
$ make run-docker
docker-compose up
Recreating completer ... done
Attaching to completer
completer    | 2019/04/27 19:37:13 config: debug: true, port: 7866, querycachettl: 20s
...

$ curl -XGET '127.0.0.1:7866/complete?term=Mos&types[]=City&types[]=airports' | jq
[
  {
	"Slug": "MOW",
	"Subtitle": "Russia",
	"Title": "Moscow"
  },
  {
	"Slug": "DME",
	"Subtitle": "Moscow",
...

$

```

## configure via ENV
```
$EDITOR docker-compose.yml
```
