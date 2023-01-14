# Setup ingress controller
1. Install ingress-nginx controller from https://kubernetes.github.io/ingress-nginx/deploy.
2. Get `EXTERNAL-IP` from `kubectl get svc ingress-nginx-controller -n ingress-nginx -ojson | jq -r '.status.loadBalancer.ingress[0].ip'`.
3. Setup an A record for DNS `yc/kubernetes/ingress/create-dns.sh --balanserip EXTERNAL-IP`.

# Setup cert-manager for Let's Encrypt
1. Instal cert-manager https://cert-manager.io/docs/installation.
2. Create `kubectl apply -f yc/kubernetes/ingress/cluster-issuer.yaml`.
