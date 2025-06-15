# Kubernetes Architecture Plan for Mycorrhizae

This document outlines a detailed plan for deploying and managing the Mycorrhizae application suite on Kubernetes, supporting two environments (`prod` and `nightly`) with cookie-based routing.

---

## 1. Cluster and Namespace Organization

1. **Kubernetes Cluster**
   - Single cluster with multiple namespaces, or separate clusters per environment (choose based on security/isolation needs).

2. **Namespaces**
   - `prod` — Production workloads.
   - `nightly` — Nightly (staging) workloads.

---

## 2. Container Images & Versioning

- Images built by CI use the monorepo commit hash as tag: `<registry>/mycorrhiza/frontend:<commit-hash>`
- Similarly for `backend-server` and `backend-ingest`.
- CI pipeline updates a `values.yaml` or Kustomize overlay with the new tag for each environment.

---

## 3. Deployments & Services

For each namespace (`prod`, `nightly`), create:

1. **Deployments**
   - `frontend` Deployment
   - `backend-server` Deployment
   - `backend-ingest` Deployment

   Key settings:
   - `replicas`: start with 2–3 replicas.
   - `image`: set via environment-specific overlay or Helm values.
   - Resource requests/limits.

2. **Services**
   - `frontend-svc` exposing port `3000`.
   - `backend-server-svc` exposing port `4041`.
   - `backend-ingest-svc` exposing port `4042`.

   Use `ClusterIP` for internal services; the Ingress will route external traffic.

---

## 4. Ingress & Cookie-Based Routing

1. **Ingress Controller**
   - Install the NGINX Ingress Controller (e.g., via Helm chart `ingress-nginx`).

2. **Ingress Resources**
   - **Prod Ingress** (in `prod` namespace)
     ```yaml
     apiVersion: networking.k8s.io/v1
     kind: Ingress
     metadata:
       name: mycorrhiza-ingress
       annotations:
         kubernetes.io/ingress.class: nginx
     spec:
       rules:
       - http:
           paths:
           - path: /
             pathType: Prefix
             backend:
               service:
                 name: frontend-svc
                 port:
                   number: 3000
           - path: /api/
             pathType: Prefix
             backend:
               service:
                 name: backend-server-svc
                 port:
                   number: 4041
           - path: /ingest/
             pathType: Prefix
             backend:
               service:
                 name: backend-ingest-svc
                 port:
                   number: 4042
     ```

   - **Nightly (Canary) Ingress** (in `prod` namespace)
     ```yaml
     apiVersion: networking.k8s.io/v1
     kind: Ingress
     metadata:
       name: mycorrhiza-canary-ingress
       annotations:
         kubernetes.io/ingress.class: nginx
         nginx.ingress.kubernetes.io/canary: "true"
         nginx.ingress.kubernetes.io/canary-by-cookie: "target_deployment_override=nightly"
     spec:
       rules:
       - http:
           paths:
           - path: /
             pathType: Prefix
             backend:
               service:
                 name: nightly-frontend-svc
                 port:
                   number: 3000
           - path: /api/
             pathType: Prefix
             backend:
               service:
                 name: nightly-backend-server-svc
                 port:
                   number: 4041
           - path: /ingest/
             pathType: Prefix
             backend:
               service:
                 name: nightly-backend-ingest-svc
                 port:
                   number: 4042
     ```

3. **Behavior**
   - Requests without `target_deployment_override` cookie route to the prod services.
   - Requests with `Cookie: target_deployment_override=nightly` go to nightly services.

---

## 5. Configuration Management

- **Helm** or **Kustomize** for templating:
  - Define charts/overlays for each component.
  - Use environment-specific values (image tags, replica counts).
- **CI/CD Integration**:
  - On commit, build images, push to registry.
  - Run `helm upgrade --install` or `kubectl apply -k` for `nightly` namespace.
  - After smoke tests, trigger deploy to `prod` namespace.

---

## 6. Rolling Updates & Rollbacks

- Leverage Kubernetes rolling update strategy (maxUnavailable/ maxSurge).
- Keep previous ReplicaSets for instant rollback.
- Monitor readiness probes and health checks.

---

## 7. Monitoring & Logging

- **Monitoring**: Prometheus + Grafana for metrics.
- **Logging**: Fluentd/Fluent Bit or EFK stack (Elasticsearch, Fluentd, Kibana).
- **Alerts**: CPU/Memory spikes, pod restarts, high error rates.

---

## 8. Next Steps

1. Prototype charts/overlays for one component.
2. Set up Ingress controller in a test cluster.
3. Validate cookie-based routing.
4. Extend to all components.
5. Integrate with CI pipeline.

---

*Document generated on 2025-06-14 by LLM.*
