# kafka-pubsub

This is a coding challenge for [goraft](https://goraft.tech/). You can find the challenge specification
in `challenge.md`.

This repository contains two executables, `cmd/publisher/publisher.go`, `cmd/subscriber/subscriber.go` that use a few
shared functions in the top level directory. These executables are built and put into Docker containers. You can find
their associated `Dockerfile`s in the top level directory.

The program can be built and reviewed by simply running `docker-compose up` in the top level directory with the correct
permissions. Here is some example truncated output:

```
publisher     | Connected to Kafka leader.
subscriber    | Connected to Kafka leader.
subscriber    | New message:
subscriber    |   Epoch in seconds: 1608129684
subscriber    |   Celsius reading: 33.35?C
subscriber    | New message:
subscriber    |   Epoch in seconds: 1608129685
subscriber    |   Celsius reading: 40.18?C
subscriber    | New message:
subscriber    |   Epoch in seconds: 1608129686
subscriber    |   Celsius reading: -7.93?C
```

Please note that the full logs for the program can be reviewed by using `docker logs publisher` or
`docker logs subscriber`. This will fix the `°` symbol not being found.

```
Connected to Kafka leader.
New message:
  Epoch in seconds: 1608129684
  Celsius reading: 33.35°C
New message:
  Epoch in seconds: 1608129685
  Celsius reading: 40.18°C
```

## Improvements

- [ ] Assign the `KAFKA_ADVERTISED_HOST_NAME` environment variable to something that will scale up with containers, if
  it doesn't already.
- [ ] Make the constants in `util.go` configurable for the executables on start up.
