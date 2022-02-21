# Webcrawler

## Description

This web crawler service provides the sitemap of the url to crawl.  
It does not follow URLs outside the root url.
The crawl request is a **GET HTTP request** with query params **url** and **depth**:
- *url*: the url to crawl.
- *depth*: the extent/level to which the webcrawler fetchs links.  
The crawl result is sent as a **JSON** payload.

For example if the crawl request is for 'http://example.com/foo', the result will look like:
```json
{
  "Path": "/foo",
  "URLs": [
    {
      "Path": "/foo/bar1",
      "URLs": [
        {
          "Path": "/foo/bar1/toto",
          "URLs": []
        },
      ]
    },
    {
      "Path": "/foor/bar2",
      "URLs": []
    },
    {
      "Path": "/foo/bar3",
      "URLs": [
        {
          "Path": "/foo/bar3/toto1",
          "URLs": []
        },
        {
          "Path": "/foo/bar3/toto2",
          "URLs": []
        },
      ]
    },
    {
      "Path": "/foo/bar3",
      "URLs": []
    }
  ]
}
```

## Usage

### Server
The main binary is the server one.  
It can be installed via `go install github.com/mo3m3n/webcrawler/cmd/server` 
```
Usage of ./webcrawler:
  -address string
        the TCP network address the webcrawler is going to listen to (default "127.0.0.1:8080")
  -path string
        the path where the webcrawser is processing crawl requests (default "/crawl")
  -log int
        the webcrawler logging level: 1=error, 2=warning, 3=info, 4=debug (default 3)
  -maxconn int
        the maximum number of concurrent requests the webcrawler can accept (default 5)
  -ratelimit int
        the maximum number of requests/second the webcrawler is allowed to send to a given website (default 1)
  -timeout int
        the number of seconds the webcrawler is going to wait for a crawl operation before interrupting it (default 300)
```

### Client
Eventhough any HTTP client can be used, this project provides a dedicated client via `go install github.com/mo3m3n/webcrawler/cmd/client`
```
webcralwer <address> <url> [depth]

  address: the TCP network address of the webcrawler.
  url: the starting url to crawl from.
  depth: the extent/level to which the webcrawler fetchs links. -1 means no limit.

```

Example:
`webcrawler http://<crawler-address>/crawl https://example.com/foo 4`

## RoadMap
- Add tests
- Add custom request header to be used by the webcrawler (example custom User-Agent)
- Honor robots.txt
- Allow the webcrawler to follow a provided list of external urls	
