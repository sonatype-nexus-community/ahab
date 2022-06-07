/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
@Library(['private-pipeline-library', 'jenkins-shared']) _

dockerizedBuildPipeline(
    pathToDockerfile: "jenkins.dockerfile",
    deployBranch: 'main',
    prepare: {
      githubStatusUpdate('pending')
    },
//    buildAndTest: {
//      sh '''
//    make all
//    '''
//    },
    // Avoid 'test', 'lint' targets for now due to docker calls having issues on Jenkins
    buildAndTest: {
      sh '''
      make deps build
    '''
    },
    vulnerabilityScan: {
      withDockerImage(env.DOCKER_IMAGE_ID, {
        withCredentials([usernamePassword(credentialsId: 'jenkins-iq',
            usernameVariable: 'IQ_USERNAME', passwordVariable: 'IQ_PASSWORD')]) {
          sh 'go list -json -deps | /tmp/tools/nancy iq --iq-application ahab --iq-stage release --iq-username $IQ_USERNAME --iq-token $IQ_PASSWORD --iq-server-url https://iq.sonatype.dev'
        }
      })
    },
    onSuccess: {
      githubStatusUpdate('success')
    },
    onFailure: {
      githubStatusUpdate('failure')
      notifyChat(currentBuild: currentBuild, env: env, room: 'community-oss-fun')
      sendEmailNotification(currentBuild, env, [], 'community-group@sonatype.com')
    }
)
