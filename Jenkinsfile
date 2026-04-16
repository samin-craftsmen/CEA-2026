pipeline {
    agent any

    environment {
        AWS_DEFAULT_REGION = "ap-south-1"
        AWS_ACCESS_KEY_ID = "YOUR_ACCESS_KEY"
        AWS_SECRET_ACCESS_KEY = "YOUR_SECRET_KEY"
    }

    stages {

        stage('Build All Lambdas') {
            steps {
                bat '''
                if not exist BACKEND\\build mkdir BACKEND\\build

                for /d %%d in (BACKEND\\lambdas\\*) do (
                    echo Building %%~nxd...

                    cd %%d

                    set GOOS=linux
                    set GOARCH=amd64

                    go build -o main

                    powershell Compress-Archive -Path main -DestinationPath %%~nxd.zip

                    move %%~nxd.zip ..\\..\\build\\
                    del main

                    cd ..\\..
                )
                '''
            }
        }

        stage('Terraform Deploy') {
            steps {
                dir('BACKEND\\terraform') {
                    bat '''
                    terraform init
                    terraform apply -auto-approve
                    '''
                }
            }
        }
    }
}
