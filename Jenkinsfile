properties([[$class: 'BuildDiscarderProperty', strategy: [$class: 'LogRotator', artifactNumToKeepStr: '2', numToKeepStr: '2']]])

node {
    def projectName = "pingpong"
    
    def gopath = pwd() + "/gopath"
    def projectDir = "${gopath}/src/github.com/WatchBeam/${projectName}"

    env.GOPATH = "${gopath}"
    env.PATH = env.PATH + ":${gopath}/bin"

    try {
        sh "mkdir -p '${projectDir}'"
        dir (projectDir) {
            stage("Checkout") {
                checkout scm
            }
            stage("Prepare") {
                sh 'go get -d ./...'
            }
            stage("go build") {
                sh "go build ./bin/${projectName}"
            }
            stage("artifacts") {
                archiveArtifacts artifacts: "${projectName}", fingerprint: false
            }
            currentBuild.result = "SUCCESS"
        }
    } catch(e) {
        currentBuild.result = "FAILURE"
        throw e
    }
}
