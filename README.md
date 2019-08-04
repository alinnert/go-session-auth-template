# Demo Auth-Server in Golang

This repo demonstrates how you can implement cookie session based Authentication in Golang. This codebase uses:

- [chi](https://github.com/go-chi/chi) (router)
- [scs](https://github.com/alexedwards/scs) (sessions)
- [Badger](https://github.com/dgraph-io/badger) (data store)

## Folder Structure

- `handlers`: Contains all handlers, and two functions to create easily uniformed responses.
- `middleware`: Contains all middlewares, and a custom type to help creating configurable middlewares.
- `models`: Contains structs that provide access to the data store.
- `server`: Contains the starting point of the application and all basic setup logic.
- `services`: Contains additional helper logic and stuff like wrapper functions.
- `values`: Contains simple values like identifiers and global keys that are needed around the app.

## Try it out

Using bash:

~~~ bash
go build
./auth-server
~~~

Using PowerShell:

~~~ powershell
go build
auth-server.exe
~~~

There's also a Makefile that does just that. If you can use the Makefile, you could instead just type:

~~~ bash
make
~~~

## Routes

### `GET /`

Test, if the server is running. The response should be:

~~~ json
{
    "data": "Demo Auth Server is running.",
    "status": "ok"
}
~~~

### `GET /db`

A route for debugging purposes. It returns all values stored in the Badger store.

### `GET /user?email={email}`

Once you added a user you can fetch its data from the data store using this route.

### `GET /profile`

If you're signed in it returns your user name. Otherwise it tells you that you're not signed in.

### `GET /public`

This should return a JSON in any case. Similar to `GET /`.

### `GET /secret`

This should only return a JSON once you're signed in. Otherwise it returns a `403` status.

### `POST /auth/signup`

Adds a user to the data store. The password gets hashed and salted using `bcrypt`.

**Request body**:

~~~ json
{
	"email": "john@doe.com",
	"password": "secret-password"
}
~~~

### `POST /auth/signin`

Signs in the user. Creates a session and adds a session cookie to the response. Currently the sessions are stored in-memory because `scs` doesn't support Badger yet.

**Request body**:

~~~ json
{
	"email": "john@doe.com",
	"password": "secret-password"
}
~~~

### `POST /auth/signout`

Deletes the session, and so the user gets signed out.