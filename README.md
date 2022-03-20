A program to handle web hooks: 
It listens on a http port for incoming web requests
(launched by these so called "web hooks") and performs actions
depending on what is specificed in its configuration file.

Clients can connect to the [NATS](https://nats.io) server and register jobs; by publishing
a `job.register` event with the following information:

```json3
{
  "name": "my job",                     // The job's name
  "path": "/some/path/{with}/{param}",  // The webhook's path (URL)
  "methods": ["GET", "POST", "PUT"],    // The allowed methods
  "reponse_headers": {                  // Optional response headers
    "x-my-super-header": "some value"
  },
  "statusCode": 203                     // Optional response status code
}
```

`juk` will register this job and when a request is made to it's `path`,
it will be processed and any client listening will be notified.

Example use case:

Let's say that each time something is "pushed" to your GitHub repository
you want to do something, for example: 

- Store information about the change
- Run some tests
- Push changes to another branch
- Etc...

You just need to "register" a job, and implement the behaviour.

An [example](example/main.go) program is provided (using the [Go](https://go.dev) programming language).

# Usage

`juk` works with a (`json`) configuration file, named `juk.json`.

*You can set the `JUK_CONFIG` environment variable to override this behaviour.*

This configuration file sets:

- The HTTP server address
- The HTTP server port
- Path of the certificate and key (Optional, only if you want to use HTTPS)
- The [NATS](http://nats.io) server host and port.

You can start `juk` by invoking

```
juk [optional config filepath]
```
# Credits

Thanks to [@dciccale](https://github.com/dciccale) for the logo (that I managed to loose) and ideas.
