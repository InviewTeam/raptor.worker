# Worker

## Enviroments
```
RABBIT - set this variable for rabbit addr
RABBIT_PORT - set this variable for rabbit port
RABBIT_LOGIN - set this variable for rabbit login
RABBIT_PASSWORD - set this variable for rabbit password
RABBIT_CHANNEL -set this variable for rabbit channel
```

## Build
```(bash)
docker build -t worker .
```

## Run
```(bash)
docker run -e RABBIT="localhost" -e RABBIT_PORT="5672" -e RABBIT_LOGIN="guest" -e RABBIT_PASSWORD="guest" -e RABBIT_CHANNEL="worker" worker
```

## K8S
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: worker
  namespace: raptor
  labels:
    app: worker
spec:
  replicas: 2
  selector:
    matchLabels:
      app: worker
  template:
    metadata:
      labels:
        app: worker
    spec:
      containers:
      - name: worker
        image: inviewteam/raptor-worker 
```