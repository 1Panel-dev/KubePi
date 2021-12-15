export function getK8sObject (kind,namespace) {
  switch (kind) {
    case "namespaces":
      return namespaceObj
    case "services":
      return serviceObj(namespace)
    case "pods":
      return podObj(namespace)
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


