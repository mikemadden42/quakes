pipeline{
    agent any
    stages{
        stage("checkout"){
            steps{
                git 'https://github.com/mikemadden42/quakes.git'
            }
        }
        stage("build"){
            steps{
                sh label: '', script: 'go build -ldflags="-s -w" -v -o quakes'
            }
        }
        stage("archive"){
            steps{
                archiveArtifacts 'quakes'
            }
        }
    }
    post{
        always{
            echo "========always========"
        }
        success{
            echo "========pipeline executed successfully ========"
        }
        failure{
            echo "========pipeline execution failed========"
        }
    }
}