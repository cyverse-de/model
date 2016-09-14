#!groovy
node {
    stage "Build"
    checkout scm

    dockerRepo = "test-${env.BUILD_TAG}"

    sh "docker build --rm -t ${dockerRepo} ."

    dockerTestRunner = "test-${env.BUILD_TAG}"
    try {
        stage "Test"
            sh "docker run --name ${dockerTestRunner} --rm ${dockerRepo}"
    } finally {
        sh returnStatus: true, script: "docker kill ${dockerTestRunner}"
        sh returnStatus: true, script: "docker rm ${dockerTestRunner}"
        sh returnStatus: true, script: "docker rmi ${dockerRepo}"
    }
}
