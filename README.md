![Imgur](http://i.imgur.com/fkDXsy4.png)

JUK is a simple webhook utility.

# Installation
You need [go](http://golang.org), and make sure it's installed
and setup [correctly](https://golang.org/doc/install#testing).

# Usage
JUK works with a (json) configuration file,
by default it's a file named juk.json.

You can set the **JUNK_CONFIG** environment variable to override this behaviour.

This configuration file sets

    - The HTTP server address

    - The HTTP server port

    - Path of the certificate and key (Optional, only if you want to use HTTPS)

    - The [NATS](http://nats.io) server host and port.

Once done you can start juk by invoking

```
juk [optional config filepath]
```

Clients can then connect to the NATS server and register jobs, by publishing
a ```job.register``` event with the following information.

```
{
  "name": "my job",                     job's name
  "path": "/some/path/{with}/{param}",  webhook path
  "methods": ["GET", "POST", "PUT"],    allowed methods
  "reponse_headers": {                  optional
    "x-my-super-header": "some value"
  },
  "statusCode": 203                    optional
}
```

JUK will register this job and when a request is made to it's ```path```
it will be processed and the job will be dispatched on the channel with
the job's ```name```, any client listening for this job is then notified.


# Credits

Thanks to [@dciccale](https://github.com/dciccale) for the logo and ideas.

