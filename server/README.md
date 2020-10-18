# divert - Backend

### API Spec

* `/api` - Index route, can use to check if the server if live
* `/api/createURL` - Create a URL, Sample Request body `{
    "original_url": "https://www.google.com",
    "shortened_url_code": "abc1234"
}`
* `/api/getAllURL` - Returns an array of JSON objects with all the long and short urls in it
* `/api/updateURL` - Update a url, request body same as that of create
* `/api/deleteURL` - Delete a url, Sample Request body `{"shortened_url_code: "abcd1234"}`
* `/api/login/` - Login route, send request in the following format `{"username": "xyz", "password": "abc"}`, will return a json web token that can be added to the x-auth-token header to make requests to all the above routes

### Environment Variables

* `PORT` - Server port number
* `MONGODB_URL` - Mongodb Server URL
* `DBNAME` - Name of the database
* `COLLECTION_NAME` - Name of the urls collection
* `AUTH_COLLECTION_NAME` - Name of auth collection
* `JWT_SIGNING_KEY` - Key to sign JWT tokens

Command to push subtree to heroku - `git subtree push --prefix server/ heroku master`