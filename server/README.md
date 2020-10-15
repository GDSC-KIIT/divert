# divert - Backend

### API Spec

* `/api` - Index route, can use to check if the server if live
* `/api/createURL` - Create a URL, Sample Request body `{
    "OriginalURL": "https://www.google.com",
    "ShortenedURLCode": "abc1234"
}`
* `/api/getAllURL` - Returns an array of JSON objects with all the long and short urls in it
* `/api/updateURL` - Update a url, request body same as that of create
* `/api/deleteURL` - Delete a url, Sample Request body `{"ShortenedURLCode: "abcd1234"}`