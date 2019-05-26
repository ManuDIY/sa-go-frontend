def  appName = 'sa-go'
def  feSvcName = "${appName}-frontend"
def  registry = "registry.linxlabs.com:5000"
def  buildTag = "${feSvcName}:${env.gitlabBranch}.${env.BUILD_NUMBER}"
def  containerTag = "${registry}/${buildTag}"

pipeline {
    agent any
    environment {
        JENKINS_AVATAR_URL="https://raw.githubusercontent.com/jenkinsci/jenkins/master/war/src/main/webapp/images/headshot.png"
    }
    stages {
        stage('Pre-Build') {
            steps {
                gitlabCommitStatus(name: 'Go env'){
                    rocketSend(
                        attachments: [[
                            title: "${env.gitlabSourceRepoName}",
                            color: 'yellow',
                            text: "In-progress :arrows_counterclockwise:\nRecent Changes - \n${getChangeString(10)}",
                            thumbUrl: '',
                            messageLink: '',
                            collapsed: false,
                            authorName: "Build job ${env.BUILD_NUMBER} started",
                            authorIcon: '',
                            authorLink: '',
                            titleLink: "${env.BUILD_URL}",
                            titleLinkDownload: '',
                            imageUrl: '',
                            audioUrl: '',
                            videoUrl: ''
                        ]],
                        channel: 'sa-project',
                        message: '',
                        avatar: "${JENKINS_AVATAR_URL}",
                        failOnError: true,
                        rawMessage: true
                    )
                    sh "env"
                }
            }
        }
        stage('Build') {
            steps {
                gitlabCommitStatus(name: 'Go build'){
                    sh '''
                    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /home/jenkins/go/bin/go build -a -installsuffix cgo -ldflags="-w -s"
                    '''
                }
            }
        }
        stage('Package and Stage') {
            steps {
                gitlabCommitStatus(name: 'Docker build and Push'){
                    script {
                        docker.withRegistry('https://registry.linxlabs.com:5000', 'docker-cred'){
                                def customImage = docker.build("${buildTag}")
                                customImage.push()
                        }
                    }
                }
            }
        }
        stage('Deploy to Kubernetes Cluster'){
            steps{
                gitlabCommitStatus(name: 'K8S Deployment'){
                        sh("sed -i.bak 's#SA_FRONTEND_CONTAINER_IMAGE#${containerTag}#' ${feSvcName}.yaml")
                        sh("/home/jenkins/kubectl --kubeconfig=/home/jenkins/k8s-cluster.yaml --namespace default apply -f ${feSvcName}.yaml")
                        }
                }
        }
        stage('Check deployment time'){
            steps{
                gitlabCommitStatus(name: 'K8S Deployment Replica status'){
                sh '''
                    count=0
                    while true
                        do
                        if [ $(/home/jenkins/kubectl --kubeconfig=/home/jenkins/k8s-cluster.yaml get deployments sa-go-frontend  --output=jsonpath={.status.availableReplicas}) -ge 1 ]
                        then
                         exit 0
                        else
                         sleep 10
                         count=`expr $count + 1`
                         if [ $count -eq 3 ]
                         then
                           echo "Unable to start Pod"
                           exit 1
                         fi
                        fi
                done'''
                }
            }
        }
        stage('App health check'){
        steps {
                gitlabCommitStatus(name: 'Application Health'){
                        timeout(time: 30, unit: 'SECONDS') {
                                retry(5) {
                                        sh '''
                                        curl -s -o /dev/null -w "%{http_code}" http://$(/home/jenkins/kubectl --kubeconfig=/home/jenkins/k8s-cluster.yaml get svc sa-go-frontend --no-headers --output=jsonpath={.status.loadBalancer.ingress[0].ip})
                                        '''
                                }
                        }
                }
            }
        }
    }
    post {
        success {
            rocketSend(
                attachments: [[
                    title: "${env.gitlabSourceRepoName}",
                    color: 'green',
                    text: "Success  :white_check_mark: \nRecent Changes - \n${getChangeString(10)}",
                    thumbUrl: '',
                    messageLink: '',
                    collapsed: false,
                    authorName: "Build job ${env.BUILD_NUMBER} completed",
                    authorIcon: '',
                    authorLink: '',
                    titleLink: "${env.BUILD_URL}",
                    titleLinkDownload: '',
                    imageUrl: '',
                    audioUrl: '',
                    videoUrl: ''
                ]],
                channel: 'sa-project',
                message: '',
                avatar: "${JENKINS_AVATAR_URL}",
                failOnError: true,
                rawMessage: true
            )
            updateGitlabCommitStatus(name: 'Pipeline', state: 'success')
        }
        failure {
            updateGitlabCommitStatus(name: 'Pipeline', state: 'failed')
            rocketSend(
                attachments: [[
                    title: "${env.gitlabSourceRepoName}",
                    color: 'red',
                    text: "Failed  :negative_squared_cross_mark: \nRecent Changes - ${getChangeString(10)}",
                    thumbUrl: '',
                    messageLink: '',
                    collapsed: false,
                    authorName: "Build job ${env.BUILD_NUMBER} failed",
                    authorIcon: '',
                    authorLink: '',
                    titleLink: "${env.BUILD_URL}",
                    titleLinkDownload: '',
                    imageUrl: '',
                    audioUrl: '',
                    videoUrl: ''
                ]],
                channel: 'sa-project',
                message: '',
                avatar: "${JENKINS_AVATAR_URL}",
                failOnError: true,
                rawMessage: true
            )
        }
    }
}
@NonCPS
def getChangeString(maxMessages) {
    MAX_MSG_LEN = 100
    COMMIT_HASH_DISPLAY_LEN = 7
    def changeString = ""

    def changeLogSets = currentBuild.changeSets


    for (int i = 0; i < changeLogSets.size(); i++) {
        def entries = changeLogSets[i].items
        for (int j = 0; j < entries.length && i + j < maxMessages; j++) {
            def entry = entries[j]
            truncated_msg = entry.msg.take(MAX_MSG_LEN)
            commitHash = entry.commitId.take(COMMIT_HASH_DISPLAY_LEN)
            changeString += "${commitHash}... - ${truncated_msg} [${entry.author}]\n"
        }
    }

    if (!changeString) {
        changeString = " There have not been changes since the last build."
    }
    return changeString
}

