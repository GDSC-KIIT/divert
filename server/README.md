# divert - Backend

### API Spec

* `/api` - Index route, can use to check if the server if live
* `/api/createURL` - Create a URL, Sample Request body `{
    "original_url": "https://www.google.com",
    "shortened_url_code": "abc1234"
}`
* `/api/getAllURL` - Returns an array of JSON objects with all the long and short urls in it
* `/api/updateURL` - Update a url, request body same as that of create
* `/api/deleteURL` - Delete a url, Sample Request body `{"shortened_url)code: "abcd1234"}`

### Environment Variables

* `PORT` - Server port number
* `MONGODB_URL` - Mongodb Server URL
* `DBNAME` - Name of the database
* `COLLECTION_NAME` - Name of the collection

Command to push subtree to heroku - `git subtree push --prefix server/ heroku master`