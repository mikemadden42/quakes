node {

stage('checkout') {
    git 'https://github.com/mikemadden42/quakes.git'
}

stage('build') {
    sh label: '', script: 'go build -ldflags="-s -w" -v -o quakes'
}

stage('archive') {
    archiveArtifacts 'quakes'
    }
}
