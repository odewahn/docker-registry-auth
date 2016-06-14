# Docker Registry Authentication in Go

This repo shows how to authenticate against the Docker registry.  I'll go (ha!) step by step, and show how to use [httpie](https://github.com/jkbrzt/httpie) as a client, and then show the corresponding go code that performs the equivalent operation.  I'll be using the great [gorequest](https://github.com/parnurzeal/gorequest) library to make all the http calls.

For other references, see:

* https://docs.docker.com/registry/spec/auth/token/
* http://www.cakesolutions.net/teamblogs/docker-registry-api-calls-as-an-authenticated-user


## Make an Initial Request to an Authenticated Service

The

```
http https://index.docker.io/v2/odewahn/myalpine/tags/list
```

Returns this:

```
HTTP/1.1 401 Unauthorized
Content-Length: 148
Content-Type: application/json; charset=utf-8
Date: Mon, 13 Jun 2016 18:52:27 GMT
Docker-Distribution-Api-Version: registry/2.0
Strict-Transport-Security: max-age=31536000
Www-Authenticate: Bearer realm="https://auth.docker.io/token",service="registry.docker.io",scope="repository:odewahn/myalpine:pull"

{
    "errors": [
        {
            "code": "UNAUTHORIZED",
            "detail": [
                {
                    "Action": "pull",
                    "Name": "odewahn/myalpine",
                    "Type": "repository"
                }
            ],
            "message": "authentication required"
        }
    ]
}
```

# Use the stuff in `Www-Authentication` header to answer the challenge

Once you have the challenge information, you have to make a request like this (note that the article does not URL encode the data, so the link breaks and curl returns a 404):

```
export PWD=<insert your password here>

http -a odewahn:$PWD https://auth.docker.io/token \
  service==registry.docker.io \
  scope==repository:odewahn/myalpine:pull
```

Returns this

```
HTTP/1.1 200 OK
Content-Length: 1494
Content-Type: application/json
Date: Mon, 13 Jun 2016 18:53:28 GMT
Strict-Transport-Security: max-age=31536000

{
    "token": "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsIng1YyI6WyJNSUlDTHpDQ0FkU2dBd0lCQWdJQkFEQUtCZ2dxaGtqT1BRUURBakJHTVVRd1FnWURWUVFERXp0Uk5Gb3pPa2RYTjBrNldGUlFSRHBJVFRSUk9rOVVWRmc2TmtGRlF6cFNUVE5ET2tGU01rTTZUMFkzTnpwQ1ZrVkJPa2xHUlVrNlExazFTekFlRncweE5qQTFNekV5TXpVNE5UZGFGdzB4TnpBMU16RXlNelU0TlRkYU1FWXhSREJDQmdOVkJBTVRPMUV6UzFRNlFqSkpNenBhUjFoT09qSlhXRTA2UTBWWFF6cFVNMHhPT2tvMlYxWTZNbGsyVHpwWlFWbEpPbGhQVTBRNlZFUlJTVG8wVWtwRE1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRVo0NkVLV3VKSXhxOThuUC9GWEU3U3VyOXlkZ3c3K2FkcndxeGlxN004VHFUa0N0dzBQZm1SS2VLdExwaXNTRFU4LzZseWZ3QUFwZWh6SHdtWmxZR2dxT0JzakNCcnpBT0JnTlZIUThCQWY4RUJBTUNCNEF3RHdZRFZSMGxCQWd3QmdZRVZSMGxBREJFQmdOVkhRNEVQUVE3VVROTFZEcENNa2t6T2xwSFdFNDZNbGRZVFRwRFJWZERPbFF6VEU0NlNqWlhWam95V1RaUE9sbEJXVWs2V0U5VFJEcFVSRkZKT2pSU1NrTXdSZ1lEVlIwakJEOHdQWUE3VVRSYU16cEhWemRKT2xoVVVFUTZTRTAwVVRwUFZGUllPalpCUlVNNlVrMHpRenBCVWpKRE9rOUdOemM2UWxaRlFUcEpSa1ZKT2tOWk5Vc3dDZ1lJS29aSXpqMEVBd0lEU1FBd1JnSWhBTzYxSWloN1FUcHNTMFFIYUNwTDFZTWNMMnZXZlNydlhHbHpSRDEwN2NRUEFpRUFtZXduelNYRHplRGxqcDc4T1NsTFFzbnROYWM5eHRyYW0xU0kxY0ZXQ2tJPSJdfQ.eyJhY2Nlc3MiOlt7InR5cGUiOiJyZXBvc2l0b3J5IiwibmFtZSI6Im9kZXdhaG4vbXlhbHBpbmUiLCJhY3Rpb25zIjpbInB1bGwiXX1dLCJhdWQiOiJyZWdpc3RyeS5kb2NrZXIuaW8iLCJleHAiOjE0NjU4NjYyMDcsImlhdCI6MTQ2NTg2NTkwNywiaXNzIjoiYXV0aC5kb2NrZXIuaW8iLCJqdGkiOiJxNDhuVnRHWkFsWkZlbm9qY1k2TyIsIm5iZiI6MTQ2NTg2NTkwNywic3ViIjoiOTU4MzExZDgtNzRjMC0xMWU0LWJlYTQtMDI0MmFjMTEwMDFiIn0.cHZsxrMtmUYIWyAh58n4Dx1kgOW4PmMvIEO6R7VHiLb1-4Z71AD6SUO4-SGK67BlRW5659LD6mX5FlitOoQHNg"
}
```

# Resubmit the original request and pass the token in the header

```
http https://index.docker.io/v2/odewahn/myalpine/tags/list 'Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsIng1YyI6WyJNSUlDTHpDQ0FkU2dBd0lCQWdJQkFEQUtCZ2dxaGtqT1BRUURBakJHTVVRd1FnWURWUVFERXp0Uk5Gb3pPa2RYTjBrNldGUlFSRHBJVFRSUk9rOVVWRmc2TmtGRlF6cFNUVE5ET2tGU01rTTZUMFkzTnpwQ1ZrVkJPa2xHUlVrNlExazFTekFlRncweE5qQTFNekV5TXpVNE5UZGFGdzB4TnpBMU16RXlNelU0TlRkYU1FWXhSREJDQmdOVkJBTVRPMUV6UzFRNlFqSkpNenBhUjFoT09qSlhXRTA2UTBWWFF6cFVNMHhPT2tvMlYxWTZNbGsyVHpwWlFWbEpPbGhQVTBRNlZFUlJTVG8wVWtwRE1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRVo0NkVLV3VKSXhxOThuUC9GWEU3U3VyOXlkZ3c3K2FkcndxeGlxN004VHFUa0N0dzBQZm1SS2VLdExwaXNTRFU4LzZseWZ3QUFwZWh6SHdtWmxZR2dxT0JzakNCcnpBT0JnTlZIUThCQWY4RUJBTUNCNEF3RHdZRFZSMGxCQWd3QmdZRVZSMGxBREJFQmdOVkhRNEVQUVE3VVROTFZEcENNa2t6T2xwSFdFNDZNbGRZVFRwRFJWZERPbFF6VEU0NlNqWlhWam95V1RaUE9sbEJXVWs2V0U5VFJEcFVSRkZKT2pSU1NrTXdSZ1lEVlIwakJEOHdQWUE3VVRSYU16cEhWemRKT2xoVVVFUTZTRTAwVVRwUFZGUllPalpCUlVNNlVrMHpRenBCVWpKRE9rOUdOemM2UWxaRlFUcEpSa1ZKT2tOWk5Vc3dDZ1lJS29aSXpqMEVBd0lEU1FBd1JnSWhBTzYxSWloN1FUcHNTMFFIYUNwTDFZTWNMMnZXZlNydlhHbHpSRDEwN2NRUEFpRUFtZXduelNYRHplRGxqcDc4T1NsTFFzbnROYWM5eHRyYW0xU0kxY0ZXQ2tJPSJdfQ.eyJhY2Nlc3MiOlt7InR5cGUiOiJyZXBvc2l0b3J5IiwibmFtZSI6Im9kZXdhaG4vbXlhbHBpbmUiLCJhY3Rpb25zIjpbInB1bGwiXX1dLCJhdWQiOiJyZWdpc3RyeS5kb2NrZXIuaW8iLCJleHAiOjE0NjU4NjYyMDcsImlhdCI6MTQ2NTg2NTkwNywiaXNzIjoiYXV0aC5kb2NrZXIuaW8iLCJqdGkiOiJxNDhuVnRHWkFsWkZlbm9qY1k2TyIsIm5iZiI6MTQ2NTg2NTkwNywic3ViIjoiOTU4MzExZDgtNzRjMC0xMWU0LWJlYTQtMDI0MmFjMTEwMDFiIn0.cHZsxrMtmUYIWyAh58n4Dx1kgOW4PmMvIEO6R7VHiLb1-4Z71AD6SUO4-SGK67BlRW5659LD6mX5FlitOoQHNg'
```

Returns this:

```
HTTP/1.1 200 OK
Content-Length: 46
Content-Type: application/json; charset=utf-8
Date: Mon, 13 Jun 2016 18:54:21 GMT
Docker-Distribution-Api-Version: registry/2.0
Strict-Transport-Security: max-age=31536000

{
    "name": "odewahn/myalpine",
    "tags": [
        "latest"
    ]
}
```
