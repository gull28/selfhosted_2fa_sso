# Self hosted 2FA single sign-on server

Single sign-on server with multiple service support.

## Features

- Add services that users can enable 2fa for.
- User registration
- Super-user registration - allows for super user to register services that are allowed to be authenticated with.
- Validate users based on active session for specific service or with a TOTP-code. (Will make a mobile app for this)

## API Docs (WIP)
*Read this if you're interested in API docs*

[https://gull28.github.io/selfhosted_2fa_sso/api-docs/](url)

## How to use

I encourage the usage of Docker because of how simple and reproducible it is.

*Regardless of the approach, you will have to create a `config.yaml` file for the project*

**Docker** (Dockerfile WIP)

To run the project using docker compose:

To build the docker image:
`docker-compose build`

To start the service:

`docker-compose up`



