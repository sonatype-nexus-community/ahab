<!--

    Copyright (c) 2019-present Sonatype, Inc.

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
CI Debug Notes
================
To validate some circleci stuff, I was able to run a “build locally” using the steps below.
The local build runs in a docker container.

  * (Once) Install circleci client (`brew install circleci`)

  * Convert the “real” config.yml into a self contained (non-workspace) config via:

        circleci config process .circleci/config.yml > .circleci/local-config.yml

  * Run a local build with the following command:
          
        circleci local execute -c .circleci/local-config.yml --job 'build'

    Typically both commands are run together:
    
        circleci config process .circleci/config.yml > .circleci/local-config.yml && circleci local execute -c .circleci/local-config.yml --job 'build'
    
    With the above command, operations that cannot occur during a local build will show an error like this:
     
      ```
      ... Error: FAILED with error not supported
      ```
    
      However, the build will proceed and can complete “successfully”, which allows you to verify scripts in your config, etc.
      
      If the build does complete successfully, you should see a happy yellow `Success!` message.

Miscellaneous
-------------

To allow your CI build to push changes back to github (e.g. release tags, etc), you need to create setup
 a github "Deploy Key" with write access. The command below will create such a key. Use an empty password.
 See: https://circleci.com/docs/2.0/add-ssh-key/#steps

    ssh-keygen -m PEM -t rsa -b 4096 -C "community-group@sonatype.com" -f <project-name>_github_rsa.key
    
Paste the public key into a new "write" enabled GitHub deploy key with Title: CircleCI Write <project name>

Be sure you check the "Allow write access" option.

    cat <project-name>_github_rsa.key.pub | pbcopy
    
In the CircleCI Web UI, under Permissions -> SSH Permissions -> Add SSH Key, enter "Hostname": github.com

Paste the private key.

    cat <project-name>_github_rsa.key | pbcopy        

As a sanity check, the private key should end with `-----END RSA PRIVATE KEY-----`.

Also update the `ssh-fingerprints:` tag in your config.yml to append the fingerprint of the write key.
