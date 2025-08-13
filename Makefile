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