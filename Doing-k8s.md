make a default run
kubectl run silicon-dawn --image=silicon-dawn:latest

save yaml
kubectl get deployment silicon-dawn -o yaml > deployment.yaml

delete old
kubectl delete deployment silicon-dawn

edit to only pull if not needed
edit deployment
 imagePullPolicy: IfNotPresent

recreate the deploy
kubectl apply -f deployment.yml

Expose the webapp
kubectl expose deployment silicon-dawn --port 8000 --target-port 8000
kubectl port-forward svc/silicon-dawn 8000

READ THIS https://kubernetes.io/docs/concepts/services-networking/service/

so editing the deployment then redeploy with 
kubectl apply -f deployment.yml



## Ingress
I am so not sure about this

set up default me ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/cloud-generic.yaml


https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/
https://kubernetes.io/docs/concepts/configuration/secret/



## How to set up my server

1. Put caddy in a docker
1. put caddy in a k8s
1. ingress caddy for access on all the ports it needs
1. Put gitea in a k8s
1. put sd in the thing
