pipeline {
    agent any

    environment {
        COMPOSE_PROJECT_NAME = 'resolve_now'

        GEMINI_API_KEY = credentials('resolve-now-gemini-api-key')
        QDRANT_API_KEY = credentials('resolve-now-qdrant-api-key')
    }

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm

                bat '''
                    git log -1 --oneline
                    echo Repository checkout completed.
                '''
            }
        }

        stage('Check Docker') {
            steps {
                bat '''
                    docker --version
                    docker compose version
                '''
            }
        }

        stage('Prepare Environment') {
            steps {
                bat '''
                    if not exist "back-end\\config" mkdir "back-end\\config"

                    (
                        echo GEMINI_API_KEY=%GEMINI_API_KEY%
                        echo QDRANT_API_KEY=%QDRANT_API_KEY%
                    ) > "back-end\\config\\.env.dev"

                    echo Environment file created.
                '''
            }
        }

        stage('Validate Compose') {
            steps {
                bat '''
                    docker compose config
                '''
            }
        }

        stage('Build') {
            steps {
                bat '''
                    docker compose build
                '''
            }
        }

        stage('Deploy') {
            steps {
                bat '''
                    docker compose down --remove-orphans

                    docker compose up -d --remove-orphans
                '''
            }
        }

        stage('Check Containers') {
            steps {
                bat '''
                    docker compose ps
                '''
            }
        }

        stage('Health Check') {
            steps {
                script {
                    sleep(time: 15, unit: 'SECONDS')
                }

                bat '''
                    echo Checking backend...

                    curl.exe --fail --silent --show-error http://localhost:8000/ || exit /b 1

                    echo.
                    echo Backend is OK.

                    echo Checking frontend...

                    curl.exe --fail --silent --show-error http://localhost:3000/ || exit /b 1

                    echo.
                    echo Frontend is OK.
                '''
            }
        }
    }

    post {

        success {
            echo '======================================'
            echo ' DEPLOYMENT SUCCESS'
            echo '======================================'
        }

        failure {
            echo '======================================'
            echo ' DEPLOYMENT FAILED'
            echo '======================================'

            echo 'Check Docker Compose manually if needed.'
        }

        always {
            echo 'Pipeline finished.'
        }
    }
}