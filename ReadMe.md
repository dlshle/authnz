# AuthNZ Platform
AuthNZ platform provides supports for Authentication and Authorization by abstracting flows and entities for common auth mechanisms and providing generalized APIs for authentication and authorization.

## Sample Config File
```
server:
  grpc: localhost:50051
database:
  host: postgres.db.host
  port: 5432
  db_name: sample_db_name
  user: authnz
  pass: ComplexPass!
```

## To Run on Docker
`docker run -d -p 50051:50051 --network auth --name authz -config=/path/to/container/config/file`

## To Run on Kubernetes
```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: authz
spec:
  replicas: 1
  selector:
    matchLabels:
      name: authz
  template:
    metadata:
      labels:
        name: authz
    spec:
      containers:
      - name: app
        image: docker.registry.com/authz/authz
        args: ["-config=path/to/config"]
        imagePullPolicy: Always
        #resources:
          #requests:
            #cpu: "0.1"
            #memory: "256Mi"
          #limits:
            #cpu: "0.2"
            #memory: "1Gi"
        ports:
          - containerPort: 50051
        env:
          - name: NODE_NAME
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
          - name: POD_NAME
            valueFrom:
              fieldRef:
                fieldPath: metadata.name
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: POD_IP
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
```