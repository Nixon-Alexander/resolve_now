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

        stage('Prepare Global Environment') {
            steps {
                bat '''
                    if not exist ".env" type nul > ".env"

                    (
                        echo GEMINI_API_KEY=%GEMINI_API_KEY%
                        echo QDRANT_API_KEY=%QDRANT_API_KEY%
                        echo PORT=8000
                        echo HOST=0.0.0.0
                        echo DB_HOST=host.docker.internal
                        echo DB_USER=user
                        echo DB_PASSWORD=user_localhost
                        echo DB_CONNECTION=tcp
                        echo DB_PORT=3306
                        echo DB_NAME=resolve_now
                        echo DB_MIGRATION_PATH=./migration/*.up.sql
                        echo QDRANT_HOST=upbeat_roentgen
                    ) > ".env"

                    echo Environment file created.
                '''
            }
        }

        stage('Prepare Dev Back End Environment') {
            steps {
                bat '''
                    if not exist "./back-end/config" mkdir "./back-end/config" 

                    (
                        echo GEMINI_API_KEY=%GEMINI_API_KEY%
                        echo QDRANT_API_KEY=%QDRANT_API_KEY%
                        echo PORT=8000
                        echo HOST=localhost
                        echo DB_HOST=localhost
                        echo DB_USER=user
                        echo DB_PASSWORD=user_localhost
                        echo DB_CONNECTION=tcp
                        echo DB_PORT=3306
                        echo DB_NAME=resolve_now
                        echo DB_MIGRATION_PATH=./migration/*.up.sql
                        echo QDRANT_HOST=upbeat_roentgen
                    ) > "./back-end/config/.env.dev"

                    echo Environment file created.
                '''
            }
        }

        stage('Prepare Prod Back End Environment') {
            steps {
                bat '''
                    if not exist "./back-end/config" mkdir "./back-end/config" 

                    (
                        echo GEMINI_API_KEY=%GEMINI_API_KEY%
                        echo QDRANT_API_KEY=%QDRANT_API_KEY%
                        echo PORT=8000
                        echo HOST=0.0.0.0
                        echo DB_HOST=host.docker.internal
                        echo DB_USER=user
                        echo DB_PASSWORD=user_localhost
                        echo DB_CONNECTION=tcp
                        echo DB_PORT=3306
                        echo DB_NAME=resolve_now
                        echo DB_MIGRATION_PATH=./migration/*.up.sql
                        echo QDRANT_HOST=qdrant
                    ) > "./back-end/config/.env.prod"

                    echo Environment file created.
                '''
            }
        }

        stage('Create network') {
            steps {
                bat '''
                    echo Checking network...

                    docker network inspect resolve_now_default >nul 2>&1

                    IF ERRORLEVEL 1 (
                        echo Network belum ada, membuat...
                        docker network create resolve_now_default
                    ) ELSE (
                        echo Network sudah ada.
                    )

                    echo Connecting container...

                    docker network inspect resolve_now_default --format="{{json .Containers}}" | findstr "upbeat_roentgen" >nul 2>&1

                    IF ERRORLEVEL 1 (
                        echo Container belum terhubung, connecting...
                        docker network connect resolve_now_default upbeat_roentgen
                    ) ELSE (
                        echo Container sudah terhubung ke network.
                    )

                    docker network ls
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

                    curl.exe --fail --silent --show-error  http://localhost:8000/api/v1/get-companies || exit /b 1

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