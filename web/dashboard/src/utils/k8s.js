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
