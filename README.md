# Docker Registry Authentication in Go

This repo shows how to authenticate against the Docker registry and retrieve the manifest for a repo.  I'll go (ha!) step by step, and show how to use [httpie](https://github.com/jkbrzt/httpie) as a client, and then show the corresponding go code that performs the equivalent operation.  I'll be using the great [gorequest](https://github.com/parnurzeal/gorequest) library to make all the http calls.

If you want to run the example, you need to set a `USER` and `PWD` environment variable.  I like to do this from the command line, like this:

```
USER=odewahn PWD=mysecret go run main.go
```

For other references and additional information, see:

* https://docs.docker.com/v1.6/registry/spec/api/
* https://docs.docker.com/registry/spec/auth/token/
* http://www.cakesolutions.net/teamblogs/docker-registry-api-calls-as-an-authenticated-user


## Make an Initial Request to an Authenticated Service

The first step is to make a request to a service that requires authentication.  In this example, I'll get the manifest information for a repo called `odewahn/myalpine:latest`

```
http https://index.docker.io/v2/odewahn/myalpine/manifests/latest
```

Here's the header and body of the response:

```
HTTP/1.1 401 Unauthorized
Content-Length: 148
Content-Type: application/json; charset=utf-8
Date: Tue, 14 Jun 2016 15:55:23 GMT
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

Note the `401 Unauthorized` return code.

# Use the stuff in `Www-Authentication` header to answer the challenge

In order to make an authorized call, we have to use the information returned in the `Www-Authenticate` header:

```
Www-Authenticate: Bearer realm="https://auth.docker.io/token",service="registry.docker.io",scope="repository:odewahn/myalpine:pull"
```

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
Date: Tue, 14 Jun 2016 15:56:53 GMT
Strict-Transport-Security: max-age=31536000

{
    "token": "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsIng1YyI6WyJNSUlDTHpDQ0FkU2dBd0lCQWdJQkFEQUtCZ2dxaGtqT1BRUURBakJHTVVRd1FnWURWUVFERXp0Uk5Gb3pPa2RYTjBrNldGUlFSRHBJVFRSUk9rOVVWRmc2TmtGRlF6cFNUVE5ET2tGU01rTTZUMFkzTnpwQ1ZrVkJPa2xHUlVrNlExazFTekFlRncweE5qQTFNekV5TXpVNE5UZGFGdzB4TnpBMU16RXlNelU0TlRkYU1FWXhSREJDQmdOVkJBTVRPMUV6UzFRNlFqSkpNenBhUjFoT09qSlhXRTA2UTBWWFF6cFVNMHhPT2tvMlYxWTZNbGsyVHpwWlFWbEpPbGhQVTBRNlZFUlJTVG8wVWtwRE1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRVo0NkVLV3VKSXhxOThuUC9GWEU3U3VyOXlkZ3c3K2FkcndxeGlxN004VHFUa0N0dzBQZm1SS2VLdExwaXNTRFU4LzZseWZ3QUFwZWh6SHdtWmxZR2dxT0JzakNCcnpBT0JnTlZIUThCQWY4RUJBTUNCNEF3RHdZRFZSMGxCQWd3QmdZRVZSMGxBREJFQmdOVkhRNEVQUVE3VVROTFZEcENNa2t6T2xwSFdFNDZNbGRZVFRwRFJWZERPbFF6VEU0NlNqWlhWam95V1RaUE9sbEJXVWs2V0U5VFJEcFVSRkZKT2pSU1NrTXdSZ1lEVlIwakJEOHdQWUE3VVRSYU16cEhWemRKT2xoVVVFUTZTRTAwVVRwUFZGUllPalpCUlVNNlVrMHpRenBCVWpKRE9rOUdOemM2UWxaRlFUcEpSa1ZKT2tOWk5Vc3dDZ1lJS29aSXpqMEVBd0lEU1FBd1JnSWhBTzYxSWloN1FUcHNTMFFIYUNwTDFZTWNMMnZXZlNydlhHbHpSRDEwN2NRUEFpRUFtZXduelNYRHplRGxqcDc4T1NsTFFzbnROYWM5eHRyYW0xU0kxY0ZXQ2tJPSJdfQ.eyJhY2Nlc3MiOlt7InR5cGUiOiJyZXBvc2l0b3J5IiwibmFtZSI6Im9kZXdhaG4vbXlhbHBpbmUiLCJhY3Rpb25zIjpbInB1bGwiXX1dLCJhdWQiOiJyZWdpc3RyeS5kb2NrZXIuaW8iLCJleHAiOjE0NjU5MjAxMTMsImlhdCI6MTQ2NTkxOTgxMywiaXNzIjoiYXV0aC5kb2NrZXIuaW8iLCJqdGkiOiJpWEJTU3k4eDlyRG5EWE9uZVBlMCIsIm5iZiI6MTQ2NTkxOTgxMywic3ViIjoiOTU4MzExZDgtNzRjMC0xMWU0LWJlYTQtMDI0MmFjMTEwMDFiIn0.b3L4IOlzs0v2asOjpVMWZBYZ1g_qP3krK08mah7De-QelLUV9KVUIOmO7tKxC0nPB6fRl0f307C1tL5rMkobRA"
}
```

# Resubmit the original request and pass the token in the header

```
http https://index.docker.io/v2/odewahn/myalpine/manifests/latest 'Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsIng1YyI6WyJNSUlDTHpDQ0FkU2dBd0lCQWdJQkFEQUtCZ2dxaGtqT1BRUURBakJHTVVRd1FnWURWUVFERXp0Uk5Gb3pPa2RYTjBrNldGUlFSRHBJVFRSUk9rOVVWRmc2TmtGRlF6cFNUVE5ET2tGU01rTTZUMFkzTnpwQ1ZrVkJPa2xHUlVrNlExazFTekFlRncweE5qQTFNekV5TXpVNE5UZGFGdzB4TnpBMU16RXlNelU0TlRkYU1FWXhSREJDQmdOVkJBTVRPMUV6UzFRNlFqSkpNenBhUjFoT09qSlhXRTA2UTBWWFF6cFVNMHhPT2tvMlYxWTZNbGsyVHpwWlFWbEpPbGhQVTBRNlZFUlJTVG8wVWtwRE1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRVo0NkVLV3VKSXhxOThuUC9GWEU3U3VyOXlkZ3c3K2FkcndxeGlxN004VHFUa0N0dzBQZm1SS2VLdExwaXNTRFU4LzZseWZ3QUFwZWh6SHdtWmxZR2dxT0JzakNCcnpBT0JnTlZIUThCQWY4RUJBTUNCNEF3RHdZRFZSMGxCQWd3QmdZRVZSMGxBREJFQmdOVkhRNEVQUVE3VVROTFZEcENNa2t6T2xwSFdFNDZNbGRZVFRwRFJWZERPbFF6VEU0NlNqWlhWam95V1RaUE9sbEJXVWs2V0U5VFJEcFVSRkZKT2pSU1NrTXdSZ1lEVlIwakJEOHdQWUE3VVRSYU16cEhWemRKT2xoVVVFUTZTRTAwVVRwUFZGUllPalpCUlVNNlVrMHpRenBCVWpKRE9rOUdOemM2UWxaRlFUcEpSa1ZKT2tOWk5Vc3dDZ1lJS29aSXpqMEVBd0lEU1FBd1JnSWhBTzYxSWloN1FUcHNTMFFIYUNwTDFZTWNMMnZXZlNydlhHbHpSRDEwN2NRUEFpRUFtZXduelNYRHplRGxqcDc4T1NsTFFzbnROYWM5eHRyYW0xU0kxY0ZXQ2tJPSJdfQ.eyJhY2Nlc3MiOlt7InR5cGUiOiJyZXBvc2l0b3J5IiwibmFtZSI6Im9kZXdhaG4vbXlhbHBpbmUiLCJhY3Rpb25zIjpbInB1bGwiXX1dLCJhdWQiOiJyZWdpc3RyeS5kb2NrZXIuaW8iLCJleHAiOjE0NjU5MjAxMTMsImlhdCI6MTQ2NTkxOTgxMywiaXNzIjoiYXV0aC5kb2NrZXIuaW8iLCJqdGkiOiJpWEJTU3k4eDlyRG5EWE9uZVBlMCIsIm5iZiI6MTQ2NTkxOTgxMywic3ViIjoiOTU4MzExZDgtNzRjMC0xMWU0LWJlYTQtMDI0MmFjMTEwMDFiIn0.b3L4IOlzs0v2asOjpVMWZBYZ1g_qP3krK08mah7De-QelLUV9KVUIOmO7tKxC0nPB6fRl0f307C1tL5rMkobRA'
```

Returns this:

```
HTTP/1.1 200 OK
Content-Length: 2007
Content-Type: application/vnd.docker.distribution.manifest.v1+prettyjws
Date: Tue, 14 Jun 2016 15:58:24 GMT
Docker-Content-Digest: sha256:dbe5a6b3f06ee68180a27bfce174c203fb58e1ab8bad0450db38a693e7b59f28
Docker-Distribution-Api-Version: registry/2.0
Etag: "sha256:dbe5a6b3f06ee68180a27bfce174c203fb58e1ab8bad0450db38a693e7b59f28"
Strict-Transport-Security: max-age=31536000

{
   "schemaVersion": 1,
   "name": "odewahn/myalpine",
   "tag": "latest",
   "architecture": "amd64",
   "fsLayers": [
      {
         "blobSum": "sha256:d0ca440e86378344053c79282fe959c9f288ef2ab031411295d87ef1250cfec3"
      }
   ],
   "history": [
      {
         "v1Compatibility": "{\"architecture\":\"amd64\",\"config\":{\"Hostname\":\"27c9668b3d5e\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":null,\"Image\":\"\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":null},\"container\":\"27c9668b3d5e3a2abeefdb725e1ff739cedda4b19eff906336298608f635b00e\",\"container_config\":{\"Hostname\":\"27c9668b3d5e\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ADD file:614a9122187935fccfa72039b9efa3ddbf371f6b029bb01e2073325f00c80b9f in /\"],\"Image\":\"\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":null},\"created\":\"2016-05-06T14:56:49.723208146Z\",\"docker_version\":\"1.9.1\",\"id\":\"e90e88b55e101f3a2752a8b784da0956c328d58eb7fdb216de4b1920bb47cee7\",\"os\":\"linux\"}"
      }
   ],
   "signatures": [
      {
         "header": {
            "jwk": {
               "crv": "P-256",
               "kid": "4MZL:Z5ZP:2RPA:Q3TD:QOHA:743L:EM2G:QY6Q:ZJCX:BSD7:CRYC:LQ6T",
               "kty": "EC",
               "x": "qmWOaxPUk7QsE5iTPdeG1e9yNE-wranvQEnWzz9FhWM",
               "y": "WeeBpjTOYnTNrfCIxtFY5qMrJNNk9C1vc5ryxbbMD_M"
            },
            "alg": "ES256"
         },
         "signature": "j2KT__uoK7wzCf38RKfhtQFaRaDO3hoo0eVHciYf6MGuUW03_Gpjg3Ks8QxyCdjIkSkaBo20Jb6Qwwcg5kiJOg",
         "protected": "eyJmb3JtYXRMZW5ndGgiOjEzNjAsImZvcm1hdFRhaWwiOiJDbjAiLCJ0aW1lIjoiMjAxNi0wNi0xNFQxNTo1ODoyNFoifQ"
      }
   ]
}
```
