FROM alpine/helm:3.10.2
COPY yc/kubernetes/ingress/helm/entrypoint.sh /usr/bin/entrypoint.sh
ENTRYPOINT ["/usr/bin/entrypoint.sh"]
