

```shell
docker build -t huiwq1990/apiserver-proxy . 
docker push huiwq1990/apiserver-proxy
```



```shell



cat<<EOF | kubectl apply -f -
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
  name: kube-proxy-service

---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-proxy-service
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
  - '*'

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-proxy-service
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kube-proxy-service
subjects:
- kind: ServiceAccount
  name: kube-proxy-service
  namespace: default

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-proxy-service
  labels:
    app: kube-proxy-service
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kube-proxy-service
  template:
    metadata:
      labels:
        app: kube-proxy-service
    spec:
      serviceAccount: kube-proxy-service
      containers:
      - name: hello-k8s
        image: huiwq1990/apiserver-proxy 
        imagePullPolicy: Always
        ports:
        - containerPort: 8080

---
apiVersion: v1
kind: Service
metadata: 
  name: kube-proxy-service
spec:
  type: ClusterIP
  selector: 
    app: kube-proxy-service
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kube-proxy-service
spec:
  rules:
  - host: kubeapi.yurt.epaas.domain
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: kube-proxy-service
            port:
              number: 80

EOF

curl -H "Host: kubeapi.yurt.epaas.domain" http://10.254.107.221/yurt/k8s/clusters/xxx/api/v1/namespaces/kube-system/pods

curl http://10.244.0.54:8080/yurt/k8s/clusters/xxx/api/v1/namespaces/kube-system/pods

curl http://127.0.0.1:8080/yurt/k8s/clusters/xxx/api/v1/namespaces/kube-system/pods


curl http://kubeapi.yurt.epaas.domain/yurt/k8s/clusters/xxx/api/v1/namespaces/kube-system/pods

```