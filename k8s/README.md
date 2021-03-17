# Kubernetes

A writeup of the deploying a Golang app with a Postgres database on Minikube.

## Minikube

```
minikube start|stop|delete

# You can ssh into the kube
minikube ssh

# or get the (public facing) ip of the kube
minikube ip
```

If K8s services are running inside Minikube, you can retrieve information with commands like:

```
# info of a service (e.g. go-observe-service)
minikube service <service-name>

# Or just get the url of the service
# where service-name is defined by metadata: name in the service yml
minikube service <service-name> --url

# e.g.
minikube service go-observe-service --url
minikube service postgres-loadbalancer  --url

```

## Elements

### Pods

Pods are the smallest deployable units that you can create and manage in Kubernetes.
A Pod is a group of one or more containers, with shared storage and network resources.

### Services

A service is basically an abstract way to expose a set of Pods to the network. 
It can have an outward facing ip, or only an internally. 
With the latter one k8s own DNS will take care of the discoverability of the service.  

### Ingres

Ingress exposes routes/endpoints from outside the cluster to services within the cluster. 
It uses name-based virtual hosting. 
Additionally, it may provide SSL and load balancing.

For Minikube, Ingress is an add on.

```
minikube addons enable ingress
```

[Nginx based ingress docs](https://kubernetes.github.io/ingress-nginx/deploy/#minikube)

## Tools

## kubectl

Applying k8s configuration:
```
# Apply a configuration to a resource by filename or stdin
kubectl apply -f k8s-web.yml

# you can ',separate' multiple files:
kubectl apply -f k8s-web.yml,k8s-web-service.yml
```

### Listing

Listing services and pods
```
kubectl get services
kubectl get pods
# for more info describe can be used instead of get
```

### Scaling
```
kubectl scale deploy go-observe --replicas=0
kubectl scale deploy go-observe --replicas=3
```

After changes in one of the yml files, the changes can be loaded in pods with
download, upscale.

Certainly not production ready, but it works fine when you're testing on a local minikube.

```
kubectl apply -f k8s-web.yml
kubectl scale deploy go-observe --replicas=0
kubectl scale deploy go-observe --replicas=1

# and check the pod id
kubectl get pods
```

### Mounting volumes

Using Minikube, two steps are needed to mount a local dir as a volume in K8s pods.

First, mount the local volume to Minikube at start, using something similar to:

```
# (!!) mind the last --mount
minikube start --mount-string="$HOME/go/src/github.com/grrrben/observe/template:/mnt/template" --mount
```

Next, use the deploy yml of the node to add the shared volume to the pod as well:
```
spec:
  ...
  template:
    ...
    spec:
      containers:
        - name: <name here>
          ...
          volumeMounts:
            - name: template
              mountPath: /app/template
      volumes:
        - name: template
          hostPath:
            path: /mnt/template
```

A bit annoying was getting multiple mounts. If you start Minikube with more than one `--mount-string` only the last 
string is mounted. And if you later try to add one more Minikube is harsh:  

`GUEST_MOUNT_CONFLICT: Sorry, docker does not allow mounts to be changed after container creation`

Mounting a level lower did solve my problem...

### Debugging

```
# general info via describe
# handy as some errors like CreateContainerError do not have log data yet
kubectl describe pods

# tailing a pod's log
# all output (println) is written to the log
# this only tails the log of a single pod,
kubectl logs <pod-name-and-id> -f=true

# logging in a pod
kubectl exec -ti <pod-name-and-id> -- /bin/bash
```

Printing all `env` vars from a pod:

```
kubectl exec <pod-name-and-id> -- env 
```

## K8s Networking

### Networking

Kubernetes has it's own DNS. 

```
$ kubectl get svc -n kube-system

NAME       TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                  AGE
kube-dns   ClusterIP   10.96.0.10   <none>        53/UDP,53/TCP,9153/TCP   2d3h

```
Whenever a service is created, k8s adds a corresponding DNS entry. 
This entry is of the form <service-name>.<namespace-name>.svc.cluster.local.

E.g., the postgres loadbalancer in our `default` namespace has a DNS record of 
`postgres-loadbalancer.default.svc.cluster.local`.

We can query this dns with `nslookup` using the service name (metadata name property) `postgres-loadbalancer` in this example. 
So for documentations sake, let's dive into it with a separate curl box. 
Remember we need to be _in_ our cluster to be able to perform the DNS queries.

```
$ kubectl run curl --image=radial/busyboxplus:curl -i --tty

# wait a bit
# and then, from the box

$ nslookup postgres-loadbalancer
Server:    10.96.0.10
Address 1: 10.96.0.10 kube-dns.kube-system.svc.cluster.local

Name:      postgres-loadbalancer
Address 1: 10.103.72.153 postgres-loadbalancer.default.svc.cluster.local

$ nslookup go-observe-service
Server:    10.96.0.10
Address 1: 10.96.0.10 kube-dns.kube-system.svc.cluster.local
  
Name:      go-observe-service
Address 1: 10.109.55.243 go-observe-service.default.svc.cluster.local

# just a reminder:
$ exit
Session ended, resume using 'kubectl attach curl -c curl -i -t' command when the pod is running
```

Relevant doc's:
- [Namespaces and DNS](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/#namespaces-and-dns)
- [Services, Load Balancing, and Networking](https://kubernetes.io/docs/concepts/services-networking/)
- [Connecting Applications with Services](https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/)

## env, configmaps and secrets

### Configmaps

Configmaps hold configuration properties in a simple key value format. 
These properties can be distributed into the pods, by referring to them in the yml's.  

```
data:
  # property-like keys; each key maps to a simple value
  db_user: "myusername"
  db_password: "mysecretpassword"
```

Example in ./k8s/config-map.yml

There is no secrecy at all with config maps. If that is needed, k8s gives you Secrets.

> ConfigMap does not provide secrecy or encryption. If the data you want to store are confidential, use a Secret rather than a ConfigMap

### env

Getting the configmaps values in the `env` from a pod is done with a configMapKeyRef:

```
  env:
    # Define the environment variable
    - name: DB_USER # Notice that the case is different here
      # from the key name in the ConfigMap.
      valueFrom:
        configMapKeyRef:
          name: configmap # The ConfigMap this value comes from.W
          key: db_user # The key to fetch.
```

If a config map is changed, the value is not in the `env` of the pod before the pod is refreshed again. 
When testing, you can scale down and up again, but this will give you downtime:

```
# first apply the new config
kubectl apply -f ./k8s/config-map.yml  

# scale down, and up. 
kubectl scale deploy go-observe --replicas=0
kubectl scale deploy go-observe --replicas=1
```

A more convenient method of refreshing all the pods is a restart rollout:

```
kubectl rollout restart deployment go-observe

# When you have a couple of pods, you can track the progress with the following command.
kubectl get pod -w
```

Showing the `env` vars from a pod:

```
kubectl exec <pod-name-and-id> -- env 
```

### Secrets

Some things have a certain amount of secrecy around it. Database passwords for example. 
K8s does of course have a way of keeping it secret. But caution should be taken and more in depth knowledge 
is needed to make your app really safe. 

> Kubernetes Secrets are, by default, stored as unencrypted base64-encoded strings.

That said, we just need secrets for our database credentials. 
So, first create base64 strings of the username and password.

```
db_username
echo -n postgres | base64
cG9zdGdyZXM=

db_password
echo -n mysecretpassword | base64
bXlzZWNyZXRwYXNzd29yZA==
```

These base64 encoded strings can be put in a secret yml:

```
data:
  # key: base64 value
  db_user: cG9zdGdyZXM=
  db_password: bXlzZWNyZXRwYXNzd29yZA==
```

Consuming the secrets is done in the yml of the resource that needs it. 

```
  env:
    # Define the environment variable
    - name: DB_USER
      valueFrom:
        secretKeyRef:
          name: secrets # The Secret this value comes from.
          key: db_user # The key to fetch.
```

todo

Relevant docs:  
- [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Creating Secrets configfiles](https://kubernetes.io/docs/tasks/configmap-secret/managing-secret-using-config-file/)


## localhost

A domain name is better than an ip. For this, you need Ingress. 
Ingress exposes routes/endpoints from outside the cluster to services within the cluster.

For Minikube, Ingress is an add on.

```
minikube addons enable ingress
```

And while you're at the terminal, take a note of the Minikube ip and add it to `/etc/hosts`. 
```
minikube ip
```

Next up, create an Ingress.yml file, and point it to your service. 
As a last step, apply your changes.

```
kubectl apply -f ingress.yml
```

The k8s docs have a chapter especially on Ingres with minikube, 
[set up Ingress on Minikube](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/).
If you use nginx, find the docs on including it in your cluster here; [nginx based ingress docs](https://kubernetes.github.io/ingress-nginx/deploy/#minikube).

### TLS

Generate a self-signed certificate and private key with:
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./cert/${KEY_FILE} -out ./cert/${CERT_FILE} -subj "/CN=${HOST}/O=${HOST}"
```

Then create the secret in the cluster via:
```
kubectl create secret tls ${CERT_NAME} --key ./cert/${KEY_FILE} --cert ./cert/${CERT_FILE}
```

[Kubernetes docs](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls)
[nginx specific docs](https://kubernetes.github.io/ingress-nginx/user-guide/tls/)

## Releases

### Rolling Releases

When a new version of a Docker container is available, you can point the image in your `.yml` to it by altering the desired version.
After applying the changes with `kubectl` the new release is rolled out.

```
kubectl apply -f web.yml
```

If you just have some minor changes in `env` vars or something similar, you can just refresh your pods.

```
kubectl rollout restart deployment go-observe

# When you have a couple of pods, you can track the progress with the following command.
kubectl get pod -w
```

### Updating a Docker container

First, add a new version to docker. E.g.

```
docker tag observe_web grrrben/go-observe:0.0.4
docker push grrrben/go-observe:0.0.4

docker tag observe-db grrrben/postgres-observe:0.0.5
docker push grrrben/postgres-observe:0.0.5
```
Next, change the version in the corresponding `.yml` file and roll out the release with an apply.

```
kubectl apply -f k8s-web.yml
```

## Read more, references

### Official
- [docs](https://kubernetes.io/docs/home/)
### Basics
- [Deploying a containerized Go app on Kubernetes](https://www.callicoder.com/deploy-containerized-go-app-kubernetes/)
- A long but insightful beginners level talk about [Docker and Kubernetes](https://youtu.be/Wf2eSG3owoA?t=10655) [Youtube Alert]

