# Golang RabbitMQ Example

This repository contains a simple example of how to use RabbitMQ with Golang. It demonstrates how to publish and consume
messages using the RabbitMQ Go client library.

## Prerequisites

- [Docker](https://www.docker.com/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Make](https://www.gnu.org/software/make/)

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/ngdangkietswe/go-rabbitmq.git
cd go-rabbitmq
```

### 2. Run project with Docker

```bash
make docker-run-prod
```

### 3. Run project with Minikube

#### 3.1. Start Minikube

```bash
minikube start
```

#### 3.2. Apply

```bash
make kube-apply
```

#### 3.3. Forward port

```bash
make kube-forward-notificationapi-port
```

### 4. Access the Application

#### 4.1. Health Check

```bash
curl http://localhost:3000/api/v1/health
```

#### 4.2. Publish a Message

```bash
curl -X POST http://localhost:3000/api/v1/notifications -H "Content-Type: application/json" -d '{"type": "email", "recipient", "ngdangkietswe@yopmail.com", "title": "Test Email", "message": "This is a test email message."}'
```


