pipeline {
    agent {
        docker { image 'rpi-chroot' }
    }

    environment { 
        PATH+EXTRA = "/usr/local/go/bin"
    }

    stages {     
        stage('Get Deps') {
            steps {
                echo 'Installing dependencies'
                sh 'go version'
                sh 'go get -d'
            }
        }
        
        stage('Test') {
            steps {
                echo 'Running test'
                sh 'go test -v'
                echo 'Running test for race conditions'
                sh 'go test -race'
                echo 'Running benchmarks'
                sh 'go test -bench=.'
                echo 'Test coverage'
                sh 'go test -cover'
            }
        }        

        stage('Build') {
            steps {
                echo 'Compiling and building'
                sh 'go build'
            }
        }
        
    }
}