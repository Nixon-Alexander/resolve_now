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

                sh '''
                    git log -1 --oneline
                    echo "Repository checkout completed."
                '''
            }
        }

        stage('Check Docker') {
            steps {
                sh '''
                    docker --version
                    docker compose version
                '''
            }
        }

        stage('Prepare Global Environment') {
            steps {
                sh '''
                    printf '%s\\n' \
                        "GEMINI_API_KEY=${GEMINI_API_KEY}" \
                        "QDRANT_API_KEY=${QDRANT_API_KEY}" \
                        "PORT=8000" \
                        "HOST=0.0.0.0" \
                        "DB_HOST=host.docker.internal" \
                        "DB_USER=user" \
                        "DB_PASSWORD=user_localhost" \
                        "DB_CONNECTION=tcp" \
                        "DB_PORT=3306" \
                        "DB_NAME=resolve_now" \
                        "DB_MIGRATION_PATH=./migration/*.up.sql" \
                        "QDRANT_HOST=qdrant" \
                        > .env

                    echo "Global environment file created."
                '''
            }
        }

        stage('Prepare Dev Back End Environment') {
            steps {
                sh '''
                    mkdir -p ./back-end/config

                    printf '%s\\n' \
                        "GEMINI_API_KEY=${GEMINI_API_KEY}" \
                        "QDRANT_API_KEY=${QDRANT_API_KEY}" \
                        "PORT=8000" \
                        "HOST=localhost" \
                        "DB_HOST=localhost" \
                        "DB_USER=user" \
                        "DB_PASSWORD=user_localhost" \
                        "DB_CONNECTION=tcp" \
                        "DB_PORT=3306" \
                        "DB_NAME=resolve_now" \
                        "DB_MIGRATION_PATH=./migration/*.up.sql" \
                        "QDRANT_HOST=qdrant" \
                        > ./back-end/config/.env.dev

                    echo "Dev environment file created."
                '''
            }
        }

        stage('Prepare Prod Back End Environment') {
            steps {
                sh '''
                    mkdir -p ./back-end/config

                    printf '%s\\n' \
                        "GEMINI_API_KEY=${GEMINI_API_KEY}" \
                        "QDRANT_API_KEY=${QDRANT_API_KEY}" \
                        "PORT=8000" \
                        "HOST=0.0.0.0" \
                        "DB_HOST=host.docker.internal" \
                        "DB_USER=user" \
                        "DB_PASSWORD=user_localhost" \
                        "DB_CONNECTION=tcp" \
                        "DB_PORT=3306" \
                        "DB_NAME=resolve_now" \
                        "DB_MIGRATION_PATH=./migration/*.up.sql" \
                        "QDRANT_HOST=qdrant" \
                        > ./back-end/config/.env.prod

                    echo "Prod environment file created."
                '''
            }
        }

        stage('Create Network') {
            steps {
                sh '''
                    echo "Checking network..."

                    if ! docker network inspect resolve_now_default > /dev/null 2>&1; then
                        echo "Network belum ada, membuat..."
                        docker network create resolve_now_default
                    else
                        echo "Network sudah ada."
                    fi

                    echo "Connecting Qdrant container..."

                    if ! docker network inspect resolve_now_default \
                        --format='{{json .Containers}}' \
                        | grep -q "qdrant"; then

                        echo "Container belum terhubung, connecting..."
                        docker network connect resolve_now_default qdrant

                    else
                        echo "Container sudah terhubung ke network."
                    fi

                    echo "Docker networks:"
                    docker network ls
                '''
            }
        }

        stage('Validate Compose') {
            steps {
                sh '''
                    docker compose config
                '''
            }
        }

        stage('Build') {
            steps {
                sh '''
                    docker compose build
                '''
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                    docker compose down --remove-orphans

                    docker compose up -d --remove-orphans
                '''
            }
        }

        stage('Check Containers') {
            steps {
                sh '''
                    docker compose ps
                '''
            }
        }

        stage('Health Check') {
            steps {
                script {
                    sleep(time: 15, unit: 'SECONDS')
                }

                sh '''
                    echo "=== WHO AM I ==="
                    whoami

                    echo "=== HOST ==="
                    hostname

                    echo "=== BACKEND CONTAINER ==="
                    docker compose ps

                    echo "=== BACKEND CURL ==="
                    curl --fail --silent --show-error \
                        http://192.168.11.134:8000/api/v1/get-companies

                    echo
                    echo "Backend is OK."

                    echo "=== FRONTEND CURL ==="
                    curl --fail --silent --show-error \
                        http://192.168.11.134:3000/

                    echo
                    echo "Frontend is OK."
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