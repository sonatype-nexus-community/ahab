<!--

    Copyright 2019-Present Sonatype Inc.

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

-->
<p align="center">
    <img src="https://github.com/sonatype-nexus-community/ahab/blob/master/docs/images/ahab.png" width="350"/>
</p>
<p align="center">
    <a href="https://circleci.com/gh/sonatype-nexus-community/ahab"><img src="https://circleci.com/gh/sonatype-nexus-community/ahab.svg?style=shield" alt="Circle CI Build Status"></img></a>
</p>
<p align="center">
    <a href="https://depshield.github.io"><img src="https://depshield.sonatype.org/badges/sonatype-nexus-community/ahab/depshield.svg" alt="DepShield Badge"></img></a>
</p>

# Ahab

`ahab` is a tool to check for vulnerabilities in your apt, apk, yum or dnf powered operating systems, powered by [Sonatype OSS Index](https://ossindex.sonatype.org/).

`ahab` currently works for images that use `apt`, `apk`, `yum` or `dnf` for package management and will do its best to auto detect which package 
manager is being used by your os.

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

Ubuntu

```
$ GOOS=linux GOARCH=amd64 go build 
$ docker build -f docker/dpkg-query-autodetect/Dockerfile . -t test
```

Fedora older (yum based)

```
$ GOOS=linux GOARCH=amd64 go build 
$ docker build -f docker/yum-autodetect/Dockerfile . -t test
```

Fedora latest (dnf based)
```
$ GOOS=linux GOARCH=amd64 go build 
$ docker build -f docker/dnf-autodetect/Dockerfile . -t test
```

Alpine

```
$ GOOS=linux GOARCH=amd64 go build 
$ docker build -f docker/apk-autodetect/Dockerfile . -t test
```

Depending on the OS, you'll see Ahab run and fail (Ubuntu and Fedora) or succeed (Alpine).

### Usage

```
 $ ahab
 ______      __                    __
/\  _  \    /\ \                  /\ \
\ \ \L\ \   \ \ \___       __     \ \ \____
 \ \  __ \   \ \  _ `\   /'__`\    \ \ '__`\
  \ \ \/\ \   \ \ \ \ \ /\ \L\.\_   \ \ \L\ \
   \ \_\ \_\   \ \_\ \_\\ \__/.\_\   \ \_,__/
    \/_/\/_/    \/_/\/_/ \/__/\/_/    \/___/
  _        _                           _    _
 /_)      /_` _  _  _ _/_     _  _    (/   /_` _ . _  _   _/  _
/_) /_/  ._/ /_// //_|/  /_/ /_//_'  (_X  /   / / /_'/ //_/ _\
    _/                   _/ /
Ahab version: development
Usage:
  ahab [flags]
  ahab [command]

Available Commands:
  chase       chase is used for auditing projects with OSS Index
  help        Help about any command
  iq          iq is used for auditing your projects with Nexus IQ Server

Flags:
  -h, --help   help for ahab

Use "ahab [command] --help" for more information about a command.
```

#### OSS Index usage

```
$ ahab chase
 ______      __                    __
/\  _  \    /\ \                  /\ \
\ \ \L\ \   \ \ \___       __     \ \ \____
 \ \  __ \   \ \  _ `\   /'__`\    \ \ '__`\
  \ \ \/\ \   \ \ \ \ \ /\ \L\.\_   \ \ \L\ \
   \ \_\ \_\   \ \_\ \_\\ \__/.\_\   \ \_,__/
    \/_/\/_/    \/_/\/_/ \/__/\/_/    \/___/
  _        _                           _    _
 /_)      /_` _  _  _ _/_     _  _    (/   /_` _ . _  _   _/  _
/_) /_/  ._/ /_// //_|/  /_/ /_//_'  (_X  /   / / /_'/ //_/ _\
    _/                   _/ /
Ahab version: development
Usage:
  ahab chase [flags]

Examples:

        dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab chase
        yum list installed | ./ahab chase
        apk info -vv | sort | ./ahab chase


Flags:
  -v, -- count          Set log level, higher is more verbose
      --clean-cache     Flag to clean the database cache for OSS Index
  -h, --help            help for chase
      --loud            Specify if you want non vulnerable packages included in your output
      --no-color        Specify if you want no color in your results
      --os string       Specify a value for the operating system type you want to scan (alpine, debian, fedora). Useful if autodetection fails and/or you want to explicitly set it.
      --output string   Specify the output type you want (json, text, csv) (default "text")
      --quiet           Quiet removes the header from being printed
      --token string    Specify your OSS Index API Token
      --user string     Specify your OSS Index Username
```

#### Exclude vulnerabilities

Sometimes you'll run into a dependency that after taking a look at, you either aren't affected by, or cannot resolve for some reason. Ahab understands, and will let you 
exclude these vulnerabilities so you can get back to a passing build:

Vulnerabilities excluded will then be silenced and not show up in the output or fail your build.

We support exclusion of vulnerability either by CVE-ID (ex: `CVE-2018-20303`) or via the OSS Index ID (ex: `a8c20c84-1f6a-472a-ba1b-3eaedb2a2a14`) as not all vulnerabilities have a CVE-ID.

##### Via CLI flag
* `./ahab --exclude-vulnerability CVE-789,bcb0c38d-0d35-44ee-b7a7-8f77183d1ae2`
* `./ahab --exclude-vulnerability CVE-789,bcb0c38d-0d35-44ee-b7a7-8f77183d1ae2`

##### Via file
By default if a file named `.ahab-ignore` exists in the same directory that ahab is run it will use it, will no other options need to be passed.

If you would like to define the path to the file you can use the following
* `./ahab --exclude-vulnerability-file=/path/to/your/exclude-file`
* `./ahab --exclude-vulnerability-file=/path/to/your/exclude-file`  

The file format requires each vulnerability that you want to exclude to be on a separate line. Comments are allowed in the file as well to help provide context when needed. See an example file below.

```
# This vulnerability is coming from package xyz, we are ok with this for now
CVN-111 
CVN-123 # Mitigated the risk of this since we only use one method in this package and the affected code doesn't matter
CVN-543
``` 

It's also possible to define expiring ignores. Meaning that if you define a date on a vulnerability ignore until that date it will be ignored and once that 
date is passed it will now be reported by ahab if its still an issue. Format to add an expiring ignore looks as follows. They can also be followed up by comments 
to provide context to as why its been ignored until that date.    

```
CVN-111 until=2021-01-01
CVN-543 until=2018-02-12 #Waiting on release from third party. Should be out before this date but gives us a little time to fix it. 
```

#### Nexus IQ Server Usage

```
$ ahab iq
 ______      __                    __
/\  _  \    /\ \                  /\ \
\ \ \L\ \   \ \ \___       __     \ \ \____
 \ \  __ \   \ \  _ `\   /'__`\    \ \ '__`\
  \ \ \/\ \   \ \ \ \ \ /\ \L\.\_   \ \ \L\ \
   \ \_\ \_\   \ \_\ \_\\ \__/.\_\   \ \_,__/
    \/_/\/_/    \/_/\/_/ \/__/\/_/    \/___/
  _        _                           _    _
 /_)      /_` _  _  _ _/_     _  _    (/   /_` _ . _  _   _/  _
/_) /_/  ._/ /_// //_|/  /_/ /_//_'  (_X  /   / / /_'/ //_/ _\
    _/                   _/ /
Ahab version: development
Usage:
  ahab iq [flags]

Examples:

        dpkg-query --show --showformat='${Package} ${Version}\n' | ./ahab iq --application testapp
        yum list installed | ./ahab iq --application testapp
        apk info -vv | sort | ./ahab iq --application testapp


Flags:
  -v, -- count                   Set log level, higher is more verbose
      --application string       Specify public application ID for request (required)
  -h, --help                     help for iq
      --host string              Specify Nexus IQ Server URL (default "http://localhost:8070")
      --max-retries int          Specify maximum number of tries to poll Nexus IQ Server (default 300)
      --os string                Specify a value for the operating system type you want to scan (alpine, debian, fedora). Useful if autodetection fails and/or you want to explicitly set it.
      --oss-index-token string   Specify your OSS Index API Token
      --oss-index-user string    Specify your OSS Index Username
      --quiet                    Quiet removes the header from being printed
      --stage string             Specify stage for application (default "develop")
      --token string             Specify Nexus IQ Token/Password for request (default "admin123")
      --user string              Specify Nexus IQ Username for request (default "admin")
```

## Why Ahab?

[Captain Ahab](https://en.wikipedia.org/wiki/Captain_Ahab) was a person hell bent on killing a white whale. 

This project is called `ahab` as like the wild captain, it will kill the creation of a Docker image if any vulnerabilities are found in your installed packages.

## Installation

At current time you have a few options:

TBD

### Build from source

```
$ export GO111MODULE=on
$ make deps
$ make test
$ make build
```

## Development

`ahab` is written using Golang 1.14, so it is best you start there.

Tests can be run like `make test`

## Contributing

We care a lot about making the world a safer place, and that's why we created `ahab`. If you as well want to
speed up the pace of software development by working on this project, jump on in! Before you start work, create
a new issue, or comment on an existing issue, to let others know you are!

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
