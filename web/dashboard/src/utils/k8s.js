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


const serviceObj = (namespace) =>  {
  return {
    apiVersion: "v1",
    kind: 'Service',
    metadata: objNsMetadata(namespace)
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
  return {
    apiVersion: "apps/v1",
    kind: kind,
    metadata: objNsMetadata(namespace),
    spec: {
      replicas:1,
      selector: selectorObj,
      template: templateObj
    }
  }
}
