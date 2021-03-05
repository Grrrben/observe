kubernetes-apply-all:
	kubectl apply -f ./k8s/config-map.yml,\
./k8s/postgres-service.yml,\
./k8s/postgres.yml,\
./k8s/volume.yml,\
./k8s/secrets.yml,\
./k8s/web-service.yml,\
./k8s/web.yml,\
./k8s/ingress.yml

kubernetes-pods:
	kubectl get pods
kubernetes-services:
	kubectl get services