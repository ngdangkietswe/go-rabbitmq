docker-run-dev:
	cd deployments/docker && docker-compose -f docker-compose.dev.yaml up -d
docker-run-prod:
	cd deployments/docker && docker-compose -f docker-compose.prod.yaml up -d

docker-stop-dev:
	cd deployments/docker && docker-compose -f docker-compose.dev.yaml down
docker-stop-prod:
	cd deployments/docker && docker-compose -f docker-compose.prod.yaml down

kube-apply:
	cd deployments/k8s && kubectl apply -f .
kube-delete:
	cd deployments/k8s && kubectl delete -f .

kube-apply-rabbitmq:
	cd deployments/k8s && kubectl apply -f rabbitmq.yaml
kube-delete-rabbitmq:
	cd deployments/k8s && kubectl delete -f rabbitmq.yaml

kube-apply-notificationapi:
	cd deployments/k8s && kubectl apply -f notificationapi.yaml
kube-delete-notificationapi:
	cd deployments/k8s && kubectl delete -f notificationapi.yaml

kube-apply-notificationworker:
	cd deployments/k8s && kubectl apply -f notificationworker.yaml
kube-delete-notificationworker:
	cd deployments/k8s && kubectl delete -f notificationworker.yaml