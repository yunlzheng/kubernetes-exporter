kubernates Exporter
===================
Export kubernates resources health status

## How To Use

1, create cluster role and service account

```
kubecyl create -f deployments/rbac-setup.yml
```

2, create k8s exporter deployment

```
kubecyl create -f deployment/k8s-exporter-deployment.yml
```

## Metrics

```
# TYPE kubernates_deployment_status gauge
kubernates_deployment_status{name="artifactory",namespace="default"} 1
# HELP kubernates_node_status status of node reported by kubernates
# TYPE kubernates_node_status gauge
kubernates_node_status{name="dev-1",namespace=""} 1
# HELP kubernates_pod_status status of pod reported by kubernates
# TYPE kubernates_pod_status gauge
kubernates_pod_status{hostIP="101.37.27.205",message="",name="etcd-dev-6",namespace="kube-system",podIP="101.37.27.205",podPhase="Running",reason=""} 1
```
