# How to setup ingress
1. kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.2.0/deploy/static/provider/cloud/deploy.yaml. See https://kubernetes.github.io/ingress-nginx/deploy/ for the latest release.
2. Get `EXTERNAL-IP` from `kubectl get svc ingress-nginx-controller -n ingress-nginx -ojson | jq -r '.status.loadBalancer.ingress[0].ip'`.
3. Setup A record for DNS `yc/kubernetes/ingress/create-dns.sh`.

# How to setup cert-manager for Let's Encrypt
1. `kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.11.0/cert-manager.yaml`
2. Create `kubectl apply -f yc/kubernetes/ingress/cluster-issuer.yaml`.
