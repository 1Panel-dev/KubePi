export function getK8sObject (kind,namespace) {
  switch (kind) {
    case "namespaces":
      return namespaceObj
    case "services":
      return serviceObj(namespace)
    case "pods":
      return podObj(namespace)
    case "deployments":
      return workloadObj(namespace,'Deployment')
    case "daemonsets":
      return workloadObj(namespace,'DaemonSet')
    case "statefulsets":
      return workloadObj(namespace,'StatefulSet')
    case "jobs":
      return workloadObj(namespace,'Job')
    case "cronjobs":
      return workloadObj(namespace,'CronJob')
    case "endpoints":
      return endPointObj(namespace)
    case "ingresses":
      return ingressObj(namespace)
    case "networkpolicies":
      return npObj(namespace)
    case "configmaps":
      return configMapObj(namespace)
    case "secrets":
      return secretObj(namespace)
    case "resourcequotas":
      return resourceQuotaObj(namespace)
    case "limitranges":
      return limitRangesObj(namespace)
    case "horizontalpodautoscalers":
      return hpaObj(namespace)
    case "poddisruptionbudgets":
      return pdbObj(namespace)
    case "persistentvolumeclaims":
      return pvcObj(namespace)
    case "persistentvolumes":
      return pvObj
    case "storageclasses":
      return scObj
    default:
      return {}
  }
}

const objMetadata = {
  name: "",
  annotations: {},
  labels: {}
}

const objNsMetadata = (namespace) => {
  return {
    name: "",
    namespace: namespace,
    annotations: {},
    labels: {}
  }
}

const namespaceObj = {
  apiVersion: "v1",
  kind: 'Namespace',
  metadata: objMetadata
}

const selectorObj = {
  matchLabels:{
    "kubepi.org/name": ""
  }
}

const templateObj = {
  metadata: {
    labels: {
      "kubepi.org/name": ""
    }
  },
  spec: {
    containers:[{
      name:"container-0",
      imagePullPolicy: "IfNotPresent",
      ports: [],
      tty: false,
      resources: {
        requests: {},
        limits: {}
      },
      securityContext: {
        capabilities: {},
        seLinuxOptions: {}
      },
      volumeMounts: [],
    }],
    restartPolicy: "Always",
    volumes: [],
    dnsConfig: {},
    affinity: {},
    nodeAffinity: {},
    securityContext:{
      seLinuxOptions: {}
    }
  }
}

const podObj = (namespace) =>  {
  return {
    apiVersion: "v1",
    kind: 'Pod',
    metadata: objNsMetadata(namespace),
    spec: {}
  }
}

const workloadObj = (namespace,kind) => {
  let apiVersion = "apps/v1"
  if (kind === 'Job') {
    apiVersion = "batch/v1"
  }
  if (kind === 'CronJob') {
    apiVersion = "batch/v1beta1"
  }
  return {
    apiVersion: apiVersion,
    kind: kind,
    metadata: objNsMetadata(namespace),
    spec: {
      replicas:1,
      selector: selectorObj,
      template: templateObj
    }
  }
}

const serviceObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: 'Service',
    metadata: objNsMetadata(namespace),
    spec: {
      type: "ClusterIP",
      ports: [{
        name: "",
        port: "",
        protocol: "TCP",
        targetPort: ""
      }],
      sessionAffinity: "None",
      sessionAffinityConfig: {
        clientIP: {}
      },
      externalIPs: [],
      clusterIP: ""
    }
  }
}

const endPointObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: 'Endpoints',
    metadata: objNsMetadata(namespace),
  }
}

const ingressObj = (namespace) => {
  return {
    apiVersion: "networking.k8s.io/v1beta1",
    kind: 'Ingress',
    metadata: objNsMetadata(namespace),
    spec: {
      rules:[{
        host: "",
        http: {
          paths:[{
            backend:{
              service: {
                name: "",
                port: {
                  number:0
                }
              }
            },
            path: "",
            pathType: ""
          }]
        }
      }],
      defaultBackend: {
        service: {
          name: "",
          port: {
            number: 0
          }
        }
      },
      tls: [{
        secretName: "",
        hosts: []
      }]
    }
  }
}

const npObj = (namespace) => {
  return {
    apiVersion: "networking.k8s.io/v1",
    kind: 'NetworkPolicy',
    metadata: objNsMetadata(namespace),
    spec: {}
  }
}

const configMapObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: 'ConfigMap',
    metadata: objNsMetadata(namespace),
    spec: {}
  }
}

const secretObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: 'Secret',
    metadata: objNsMetadata(namespace),
    data: {},
    type: "Opaque"
  }
}

const resourceQuotaObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: 'ResourceQuota',
    metadata: objNsMetadata(namespace),
    spec: {
      hard: {}
    }
  }
}

const limitRangesObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: 'LimitRange',
    metadata: objNsMetadata(namespace),
    spec: {
      limits: []
    }
  }
}

const hpaObj = (namespace) => {
  return {
    apiVersion: "autoscaling/v2beta2",
    kind: 'HorizontalPodAutoscaler',
    metadata: objNsMetadata(namespace),
    spec: {
      metrics: [{
        type:"Resource",
        resource: {
          name: "cpu",
          target: {
            type: "Utilization",
            averageUtilization: 80
          }
        }
      }]
    }
  }
}

const pdbObj = (namespace) => {
  return {
    apiVersion: "policy/v1beta1",
    kind: 'PodDisruptionBudget',
    metadata: objNsMetadata(namespace),
    spec: {}
  }
}

const pvcObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: 'PersistentVolumeClaim',
    metadata: objNsMetadata(namespace),
    spec: {
      accessModes: [],
      resources: {
        requests: {
          storage: "1Gi"
        }
      }
    }
  }
}

const pvObj = {
  apiVersion: "v1",
  kind: 'PersistentVolume',
  metadata: objMetadata,
  spec: {
    local: {
      path: ""
    },
    capacity: {
      storage: "1Gi"
    },
    nodeAffinity: {},
    accessModes: []
  }
}

const scObj = {
  apiVersion: "storage.k8s.io/v1",
  kind: 'StorageClass',
  metadata: objMetadata,
  provisioner: "",
  parameters: {}
}
