# icii

[![Go Report Card](https://goreportcard.com/badge/broadcastle.co/code/icii)](https://goreportcard.com/report/broadcastle.co/code/icii)
[![GoDoc](https://godoc.org/broadcastle.co/code/icii?status.svg)](https://godoc.org/broadcastle.co/code/icii)

Stream MP3 files to icecast with go.

## API Paths

### V1

> /api/v1/

```
/user
    POST    |   Create a user account.
/user/login
    POST    |   Login and receive a token.
/user/edit
    POST    |   Update user account information.
    GET     |   Get user account information.
    DELETE  |   Delete the user account.
/station
    POST    |   Create a station.
/station/:id
    POST    |   Update a station.
    GET     |   Get station information.
    DELETE  |   Delete the station.
/station/:id/track
    POST    |   Upload a track.
/station/:id/track/:id
    POST    |   Update track information.
    GET     |   Get track information.
    DELETE  |   Delete the track.
```

