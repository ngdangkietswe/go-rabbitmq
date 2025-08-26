docker-run-rabbitmq:
	cd deployments/docker && docker-compose -f docker-compose.rabbitmq.yaml up -d
docker-run-dev:
	cd deployments/docker && docker-compose -f docker-compose.dev.yaml up -d
docker-run-prod:
	cd deployments/docker && docker-compose -f docker-compose.prod.yaml up -d

docker-stop-rabbitmq:
	cd deployments/docker && docker-compose -f docker-compose.rabbitmq.yaml down
docker-stop-dev:
	cd deployments/docker && docker-compose -f docker-compose.dev.yaml down
docker-stop-prod:
	cd deployments/docker && docker-compose -f docker-compose.prod.yaml down

kube-apply:
	cd deployments/k8s && kubectl apply -f .
kube-delete:
	cd deployments/k8s && kubectl delete -f .

kube-apply-rabbitmq:
	cd deployments/k8s && kubectl apply -f docker-compose.rabbitmq.yaml
kube-delete-rabbitmq:
	cd deployments/k8s && kubectl delete -f docker-compose.rabbitmq.yaml

kube-apply-notificationapi:
	cd deployments/k8s && kubectl apply -f notificationapi.yaml
kube-delete-notificationapi:
	cd deployments/k8s && kubectl delete -f notificationapi.yaml

kube-apply-notificationworker:
	cd deployments/k8s && kubectl apply -f notificationworker.yaml
kube-delete-notificationworker:
	cd deployments/k8s && kubectl delete -f notificationworker.yaml

kube-forward-notificationapi-port:
	cd deployments/k8s && kubectl port-forward svc/notificationapi 3000:3000

kube-apply-dev:
	cd deployments/k8s && kubectl apply -k overlays/dev

kube-apply-prod:
	cd deployments/k8s && kubectl apply -k overlays/prod

kube-delete-dev:
	cd deployments/k8s && kubectl delete -k overlays/dev

kube-delete-prod:
	cd deployments/k8s && kubectl delete -k overlays/prod

terraform-apply:
	cd scripts && bash terraform.sh

swag-init:
	swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal