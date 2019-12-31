# Golang Keycloak client

This repo contains the code that can be used to authenticate microservices implemented in Golang against Keycloak (or any OIDC provider).

## About the app

This app implements the logic to connect to an OIDC provider (Keycloak in this case) and fetch the authenticated user's information (name, email, username, roles, scopes, etc.). It also implements features like verifying an already acquired access token against the server and refreshing the token before / after the access token expires as long as the refresh token is valid.

### How is the app built?

- The backend http server routes have been implemented in Golang using the http module
- The UI is a static page that is rendered using Golang templates
- Styling is done using the Materialize CSS library
- jQuery is being used for the Ajax calls to the backend server
- The app uses Go modules for package management

## What can the app do?

- Authenticate against the OIDC provider to get an ID Token
- Exchange the ID Token for an access and a refresh token
- Verify the tokens received against the OIDC provider
- Refresh the access token using the refresh token
- Spins up a simple static html page that shows the user data, tokens and authentication codes obtained from the OIDC provider

## Setup and running the app

For the app to work, create a new client (`golang-client` is the default name, but configurable) in the OIDC provider page and add `openid` to the scopes.

Change the `keyCloakServerURL` and other parameters accordingly in the `main.go` file.

Start the app server using
```
$ go run main.go
```
Since the app uses go modules package management, all of the dependencies are installed implicitly.
The UI would come up at `0.0.0.0:3000`

## Contributing

I'm still a noob, trying to get a hang of how things work in Golang. Hence, I'm sure there are a lot of things that can be done better / right.

Please feel free to create an issue in case you find anything like that and I would be more than happy to address that.
