Docker Actions:

Build the image:
    docker build -t vidyablr/mcluster-valcontroller:0.1.0 .

Run the image:
    docker run --rm -ti vidyablr/mcluster-valcontroller:0.1.0

Push the image into docker:
    docker push vidyablr/mcluster-valcontroller:0.1.0 

Create a k8s deployment:
    kubectl create deployment vcontroller --image vidyablr/mcluster-valcontroller:0.1.0 --dry-run=client -oyaml > manifests/deploy.yaml

Steps to create the secret with .crt and key files:
mkdir manifests; cd manifests
mkdir certs ; cd certs
openssl req -new -x509 -subj "/CN=vcontroller.default.svc" -addext "subjectAltName = DNS:vcontroller.default.svc" -nodes -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365


Create secret:
kubectl create secret generic certs --from-file tls.crt --from-file tls.key --dry-run=client -oyaml > secret.yaml

Mount the secret on the deploy.yaml

Deploy:
kubectl create -f manifests/certs/secret.yaml
kubectl create -f manifests/deploy.yaml

Create a service:
kubectl expose deployment vcontroller --port 443 --dry-run=client -oyaml > manifests/service.yaml

Target port: 
8443

Check the pods:
kubectl get pods
kubectl describe pod <podname>
