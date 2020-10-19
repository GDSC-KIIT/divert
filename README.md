# divert

The internal URL shortening service for DSC KIIT. 

## Get Started

Install the CLI tool to login and create short urls (only for members of DSC KIIT). The package is available on npm.

```
npm i -g @dsckiit/divert
```

This will install the `divert` command, just run `divert` in your terminal to start the cli tool.

## Tech and Design

The long and short urls are stored in the database. The backend fetches the database every 3 minutes
and updates a local in memory hashmap that maps the short urls to long urls. The redirect requests are directly served 
by accessing this in memory hashmap. We need that sweet O(1) access.

The process to fetch from the database runs on another thread (implemented using goroutines). 
`sync.RWMutex` is used to prevent data races on the in memory hashmap.

The backend service is written in Go. The CLI tool is written in Typescript.

**Backend Dependencies**

* github.com/dgrijalva/jwt-go v3.2.0
* github.com/gorilla/mux v1.8.0
* go.mongodb.org/mongo-driver v1.4.2

**CLI Dependencies**

* axios
* boxen
* chalk
* configstore
* inquirer
* ora

## License

Copyright (c) **DSC KIIT**. All rights reserved. Licensed under the MIT License

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/uses-js.svg)](https://forthebadge.com)