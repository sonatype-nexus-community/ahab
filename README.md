<p align="center">
    <img src="https://github.com/sonatype-nexus-community/ahab/blob/master/docs/images/ahab.png" width="350"/>
</p>
<p align="center">
    <a href="https://travis-ci.org/sonatype-nexus-community/ahab"><img src="https://travis-ci.org/sonatype-nexus-community/ahab.svg?branch=master" alt="Build Status"></img></a>
</p>
<p align="center">
    <a href="https://depshield.github.io"><img src="https://depshield.sonatype.org/badges/sonatype-nexus-community/ahab/depshield.svg" alt="DepShield Badge"></img></a>
</p>

# Ahab

`ahab` is a tool to check for vulnerabilities in your apt or yum powered operating systems, powered by [Sonatype OSS Index](https://ossindex.sonatype.org/).

To use `ahab`, assuming you have a built version of it:

* `dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase`
* `yum list installed | ./ahab chase --os fedora`

`ahab` currently works for images that use `apt` or `yum` for package management.

## Why is this useful?

Well, we'd hope it is easy enough to see why, but what you can do with `ahab` is inject a command similar to the following in your `Dockerfile`:

```
    RUN dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase
```

Since `ahab` will exit with a non zero code if vulnerabilities are found, you can use `ahab` to prevent images with vulnerabilities from being built, serving as a gate in your CI/CD process. `ahab` does not replace checking your own applications for vulnerable dependencies, etc..., but as the container has become more and more important to how an application eventually ends up in Production, checking that base image itself is critical as well.

A suggested setup would be to have a base image similar to:

```
FROM ubuntu:latest

RUN apt-get update && apt-get install pip

RUN ./script_to_install_ahab.sh

RUN dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase
```

Using this base image, you'd install all the packages necessary to run your application, and check it as a last step with `ahab` to ensure you aren't using anything vulnerable. From here, you'd use this base image to import your application, build it, etc... as you normally would, knowing you started from a clean base.

### See it work in Docker!

In this repo we have a Dockerfile that will copy in `ahab`, and run it on Ubuntu, to illustrate a failing Docker build.

To run this test:

```
$ GOOS=linux GOARCH=amd64 go build 
$ docker build . -t test
```

You should see `ahab` run and fail the Docker build, due to some vulnerabilities in the base os packages (Ubuntu in this case)!

## Why Ahab?

[Captain Ahab](https://en.wikipedia.org/wiki/Captain_Ahab) was a person hell bent on killing a white whale. 

This project is called `ahab` as like the wild captain, it will kill the creation of a Docker image if any vulnerabilities are found in your installed packages.

## Installation

At current time you have a few options:

TBD

### Build from source

TBD

```
$ export GO111MODULE=on
$ go test ./...
$ go build
```

### Download release binary

TBD

## Development

`ahab` is written using Golang 1.12, so it is best you start there.

Tests can be run like `go test ./... -v`

## Contributing

We care a lot about making the world a safer place, and that's why we created `ahab`. If you as well want to
speed up the pace of software development by working on this project, jump on in! Before you start work, create
a new issue, or comment on an existing issue, to let others know you are!

## Acknowledgements

TBD

## The Fine Print

It is worth noting that this is **NOT SUPPORTED** by Sonatype, and is a contribution of ours
to the open source community (read: you!)

Remember:

* Use this contribution at the risk tolerance that you have
* Do NOT file Sonatype support tickets related to `ahab` support in regard to this project
* DO file issues here on GitHub, so that the community can pitch in

Phew, that was easier than I thought. Last but not least of all:

Have fun creating and using `ahab` and the [Sonatype OSS Index](https://ossindex.sonatype.org/), we are glad to have you here!

## Getting help

Looking to contribute to our code but need some help? There's a few ways to get information:

* Chat with us on [Gitter](https://gitter.im/sonatype/nexus-developers)
