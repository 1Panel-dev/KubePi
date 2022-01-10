export function getK8sObject(kind, namespace, version) {
  switch (kind) {
    case "namespaces":
      return namespaceObj;
    case "services":
      return serviceObj(namespace);
    case "pods":
      return podObj(namespace);
    case "deployments":
      return deploymentObj(namespace);
    case "daemonsets":
      return daemonsetObj(namespace);
    case "statefulsets":
      return statefulsetObj(namespace);
    case "jobs":
      return jobObj(namespace);
    case "cronjobs":
      return cronjobObj(namespace);
    case "endpoints":
      return endPointObj(namespace);
    case "ingresses":
      return ingressObj(namespace, version);
    case "networkpolicies":
      return npObj(namespace);
    case "configmaps":
      return configMapObj(namespace);
    case "secrets":
      return secretObj(namespace);
    case "resourcequotas":
      return resourceQuotaObj(namespace);
    case "limitranges":
      return limitRangesObj(namespace);
    case "horizontalpodautoscalers":
      return hpaObj(namespace);
    case "poddisruptionbudgets":
      return pdbObj(namespace);
    case "persistentvolumeclaims":
      return pvcObj(namespace);
    case "persistentvolumes":
      return pvObj;
    case "storageclasses":
      return scObj;
    case "serviceaccounts":
      return saObj(namespace);
    case "roles":
      return roleObj(namespace);
    case "rolebindings":
      return roleBindingObj(namespace);
    case "clusterroles":
      return clusterRoleObj;
    case "clusterrolebindings":
      return clusterRoleBindingObj;
    case "podsecuritypolicies":
      return pspObj;
    default:
      return {};
  }
}

const namespaceObj = {
  apiVersion: "v1",
  kind: "Namespace",
  metadata: {
    name: "test1",
    annotations: {},
    labels: {},
  },
};

const podObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: "Pod",
    metadata: {
      name: "static-web",
      namespace: namespace,
      labels: {
        role: "myrole",
      },
    },
    spec: {
      containers: [
        {
          name: "web",
          image: "nginx",
          ports: [
            {
              name: "web",
              containerPort: 80,
              protocol: "TCP",
            },
          ],
        },
      ],
    },
  };
};

const deploymentObj = (namespace) => {
  return {
    apiVersion: "apps/v1",
    kind: "Deployment",
    metadata: {
      name: "nginx-deployment",
      namespace: namespace,
      labels: {
        app: "nginx",
      },
    },
    spec: {
      replicas: 3,
      selector: {
        matchLabels: {
          app: "nginx",
        },
      },
      template: {
        metadata: {
          labels: {
            app: "nginx",
          },
        },
        spec: {
          containers: [
            {
              name: "nginx",
              image: "nginx:1.14.2",
              ports: [
                {
                  containerPort: 80,
                },
              ],
            },
          ],
        },
      },
    },
  };
};

const cronjobObj = (namespace) => {
  return {
    apiVersion: "batch/v1beta1",
    kind: "CronJob",
    metadata: {
      name: "hello",
      namespace: namespace,
    },
    spec: {
      schedule: "*/1 * * * *",
      jobTemplate: {
        spec: {
          template: {
            spec: {
              containers: [
                {
                  name: "hello",
                  image: "busybox",
                  imagePullPolicy: "IfNotPresent",
                  command: [
                    "/bin/sh",
                    "-c",
                    "date; echo Hello from the Kubernetes cluster",
                  ],
                },
              ],
              restartPolicy: "OnFailure",
            },
          },
        },
      },
    },
  };
};

const statefulsetObj = (namespace) => {
  return {
    apiVersion: "apps/v1",
    kind: "StatefulSet",
    metadata: {
      name: "web",
      namespace: namespace,
    },
    spec: {
      selector: {
        matchLabels: {
          app: "nginx",
        },
      },
      serviceName: "nginx",
      replicas: 3,
      template: {
        metadata: {
          labels: {
            app: "nginx",
          },
        },
        spec: {
          terminationGracePeriodSeconds: 10,
          containers: [
            {
              name: "nginx",
              image: "k8s.gcr.io/nginx-slim:0.8",
              ports: [
                {
                  containerPort: 80,
                  name: "web",
                },
              ],
              volumeMounts: [
                {
                  name: "www",
                  mountPath: "/usr/share/nginx/html",
                },
              ],
            },
          ],
        },
      },
      volumeClaimTemplates: [
        {
          metadata: {
            name: "www",
          },
          spec: {
            accessModes: ["ReadWriteOnce"],
            storageClassName: "my-storage-class",
            resources: {
              requests: {
                storage: "1Gi",
              },
            },
          },
        },
      ],
    },
  };
};

const daemonsetObj = (namespace) => {
  return {
    apiVersion: "apps/v1",
    kind: "DaemonSet",
    metadata: {
      name: "fluentd-elasticsearch",
      namespace: namespace,
      labels: {
        "k8s-app": "fluentd-logging",
      },
    },
    spec: {
      selector: {
        matchLabels: {
          name: "fluentd-elasticsearch",
        },
      },
      template: {
        metadata: {
          labels: {
            name: "fluentd-elasticsearch",
          },
        },
        spec: {
          tolerations: [
            {
              key: "node-role.kubernetes.io/master",
              effect: "NoSchedule",
            },
          ],
          containers: [
            {
              name: "fluentd-elasticsearch",
              image: "quay.io/fluentd_elasticsearch/fluentd:v2.5.2",
              resources: {
                limits: {
                  memory: "200Mi",
                },
                requests: {
                  cpu: "100m",
                  memory: "200Mi",
                },
              },
              volumeMounts: [
                {
                  name: "varlog",
                  mountPath: "/var/log",
                },
                {
                  name: "varlibdockercontainers",
                  mountPath: "/var/lib/docker/containers",
                  readOnly: true,
                },
              ],
            },
          ],
          terminationGracePeriodSeconds: 30,
          volumes: [
            {
              name: "varlog",
              hostPath: {
                path: "/var/log",
              },
            },
            {
              name: "varlibdockercontainers",
              hostPath: {
                path: "/var/lib/docker/containers",
              },
            },
          ],
        },
      },
    },
  };
};

const jobObj = (namespace) => {
  return {
    apiVersion: "batch/v1",
    kind: "Job",
    metadata: {
      name: "pi",
      namespace: namespace,
    },
    spec: {
      template: {
        spec: {
          containers: [
            {
              name: "pi",
              image: "perl",
              command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"],
            },
          ],
          restartPolicy: "Never",
        },
      },
      backoffLimit: 4,
    },
  };
};

const serviceObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: "Service",
    metadata: {
      name: "my-service",
      namespace: namespace,
    },
    spec: {
      selector: {
        app: "MyApp",
      },
      ports: [
        {
          protocol: "TCP",
          port: 80,
          targetPort: 9376,
        },
      ],
    },
  };
};

const endPointObj = (namespace) => {
  return {
    kind: "Endpoints",
    apiVersion: "v1",
    metadata: {
      name: "elasticsearch-1",
      namespace: namespace,
      annotations: {},
      labels: {},
    },
    subsets: [
      {
        addresses: [
          {
            ip: "192.168.11.13",
          },
        ],
        ports: [
          {
            port: 9200,
          },
        ],
      },
    ],
  };
};

const ingressObj = (namespace, version) => {
  let apiVersion =
    version.indexOf("v1.20.") !== -1
      ? "networking.k8s.io/v1"
      : "networking.k8s.io/v1beta1";
  return {
    apiVersion: apiVersion,
    kind: "Ingress",
    metadata: {
      name: "minimal-ingress",
      namespace: namespace,
      annotations: {
        "nginx.ingress.kubernetes.io/rewrite-target": "/",
      },
    },
    spec: {
      rules: [
        {
          http: {
            paths: [
              {
                path: "/testpath",
                pathType: "Prefix",
                backend: {
                  service: {
                    name: "test",
                    port: {
                      number: 80,
                    },
                  },
                },
              },
            ],
          },
        },
      ],
    },
  };
};

const npObj = (namespace) => {
  return {
    apiVersion: "networking.k8s.io/v1",
    kind: "NetworkPolicy",
    metadata: {
      name: "test-network-policy",
      namespace: namespace,
    },
    spec: {
      podSelector: {
        matchLabels: {
          role: "db",
        },
      },
      policyTypes: ["Ingress", "Egress"],
      ingress: [
        {
          from: [
            {
              ipBlock: {
                cidr: "172.17.0.0/16",
                except: ["172.17.1.0/24"],
              },
            },
            {
              namespaceSelector: {
                matchLabels: {
                  project: "myproject",
                },
              },
            },
            {
              podSelector: {
                matchLabels: {
                  role: "frontend",
                },
              },
            },
          ],
          ports: [
            {
              protocol: "TCP",
              port: 6379,
            },
          ],
        },
      ],
      egress: [
        {
          to: [
            {
              ipBlock: {
                cidr: "10.0.0.0/24",
              },
            },
          ],
          ports: [
            {
              protocol: "TCP",
              port: 5978,
            },
          ],
        },
      ],
    },
  };
};

const configMapObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: "ConfigMap",
    metadata: {
      name: "game-demo",
      namespace: namespace,
    },
    data: {
      player_initial_lives: "3",
      ui_properties_file_name: "user-interface.properties",
      "game.properties":
        "enemy.types=aliens,monsters\nplayer.maximum-lives=5\n",
      "user-interface.properties":
        "color.good=purple\ncolor.bad=yellow\nallow.textmode=true   \n",
    },
  };
};

const secretObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: "Secret",
    metadata: {
      name: "secret-basic-auth",
      namespace: namespace,
    },
    type: "kubernetes.io/basic-auth",
    stringData: {
      username: "admin",
      password: "t0p-Secret",
    },
  };
};

const resourceQuotaObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: "ResourceQuota",
    metadata: {
      name: "pods-high",
      namespace: namespace,
    },
    spec: {
      hard: {
        cpu: "1000",
        memory: "200Gi",
        pods: "10",
      },
      scopeSelector: {
        matchExpressions: [
          {
            operator: "In",
            scopeName: "PriorityClass",
            values: ["high"],
          },
        ],
      },
    },
  };
};

const limitRangesObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: "LimitRange",
    metadata: {
      name: "cpu-min-max-demo-lr",
      namespace: namespace,
    },
    spec: {
      limits: [
        {
          max: {
            cpu: "800m",
          },
          min: {
            cpu: "200m",
          },
          type: "Container",
        },
      ],
    },
  };
};

const hpaObj = (namespace) => {
  return {
    apiVersion: "autoscaling/v2",
    kind: "HorizontalPodAutoscaler",
    metadata: {
      name: "php-apache",
      namespace: namespace,
    },
    spec: {
      scaleTargetRef: {
        apiVersion: "apps/v1",
        kind: "Deployment",
        name: "php-apache",
      },
      minReplicas: 1,
      maxReplicas: 10,
      metrics: [
        {
          type: "Resource",
          resource: {
            name: "cpu",
            target: {
              type: "Utilization",
              averageUtilization: 50,
            },
          },
        },
      ],
    },
  };
};

const pdbObj = (namespace) => {
  return {
    apiVersion: "policy/v1beta1",
    kind: "PodDisruptionBudget",
    metadata: {
      name: "zk-pdb",
      namespace: namespace,
    },
    spec: {
      minAvailable: 2,
      selector: {
        matchLabels: {
          app: "zookeeper",
        },
      },
    },
  };
};

const pvcObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: "PersistentVolumeClaim",
    metadata: {
      name: "myclaim",
      namespace: namespace,
    },
    spec: {
      accessModes: ["ReadWriteOnce"],
      volumeMode: "Filesystem",
      resources: {
        requests: {
          storage: "8Gi",
        },
      },
      storageClassName: "slow",
      selector: {
        matchLabels: {
          release: "stable",
        },
        matchExpressions: [
          {
            key: "environment",
            operator: "In",
            values: ["dev"],
          },
        ],
      },
    },
  };
};

const pvObj = {
  apiVersion: "v1",
  kind: "PersistentVolume",
  metadata: {
    name: "pv0003",
  },
  spec: {
    capacity: {
      storage: "5Gi",
    },
    volumeMode: "Filesystem",
    accessModes: ["ReadWriteOnce"],
    persistentVolumeReclaimPolicy: "Recycle",
    storageClassName: "slow",
    mountOptions: ["hard", "nfsvers=4.1"],
    nfs: {
      path: "/tmp",
      server: "172.17.0.2",
    },
  },
};

const scObj = {
  apiVersion: "storage.k8s.io/v1",
  kind: "StorageClass",
  metadata: {
    name: "standard",
  },
  provisioner: "kubernetes.io/aws-ebs",
  parameters: {
    type: "gp2",
  },
  reclaimPolicy: "Retain",
  allowVolumeExpansion: true,
  mountOptions: ["debug"],
  volumeBindingMode: "Immediate",
};

const saObj = (namespace) => {
  return {
    apiVersion: "v1",
    kind: "ServiceAccount",
    metadata: {
      name: "build-robot",
      namespace: namespace,
    },
    secrets: [
      {
        name: "build-robot-token-bvbk5",
      },
    ],
  };
};

const roleObj = (namespace) => {
  return {
    apiVersion: "rbac.authorization.k8s.io/v1",
    kind: "Role",
    metadata: {
      namespace: namespace,
      name: "pod-reader",
    },
    rules: [
      {
        apiGroups: [""],
        resources: ["pods"],
        verbs: ["get", "watch", "list"],
      },
    ],
  };
};

const roleBindingObj = (namespace) => {
  return {
    apiVersion: "rbac.authorization.k8s.io/v1",
    kind: "RoleBinding",
    metadata: {
      name: "read-pods",
      namespace: namespace,
    },
    subjects: [
      {
        kind: "User",
        name: "jane",
        apiGroup: "rbac.authorization.k8s.io",
      },
    ],
    roleRef: {
      kind: "Role",
      name: "pod-reader",
      apiGroup: "rbac.authorization.k8s.io",
    },
  };
};

const clusterRoleObj = {
  kind: "ClusterRole",
  apiVersion: "rbac.authorization.k8s.io/v1",
  metadata: {
    name: "read-all-clusterrole",
  },
  rules: [
    {
      nonResourceURLs: ["/metrics"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: [""],
      resources: [
        "bindings",
        "componentstatuses",
        "configmaps",
        "endpoints",
        "events",
        "limitranges",
        "namespaces",
        "namespaces/finalize",
        "namespaces/status",
        "nodes",
        "nodes/proxy",
        "nodes/status",
        "persistentvolumeclaims",
        "persistentvolumeclaims/status",
        "persistentvolumes",
        "persistentvolumes/status",
        "pods",
        "pods/attach",
        "pods/binding",
        "pods/eviction",
        "pods/exec",
        "pods/log",
        "pods/proxy",
        "pods/status",
        "podtemplates",
        "replicationcontrollers",
        "replicationcontrollers/scale",
        "replicationcontrollers/status",
        "resourcequotas",
        "resourcequotas/status",
        "serviceaccounts",
        "services",
        "services/proxy",
        "services/status",
      ],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["apps"],
      resources: [
        "controllerrevisions",
        "daemonsets",
        "daemonsets/status",
        "deployments",
        "deployments/scale",
        "deployments/status",
        "replicasets",
        "replicasets/scale",
        "replicasets/status",
        "statefulsets",
        "statefulsets/scale",
        "statefulsets/status",
      ],
      verbs: ["list", "get", "watch"],
    },
    {
      apiGroups: ["batch"],
      resources: ["jobs", "jobs/status"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["autoscaling"],
      resources: [
        "horizontalpodautoscalers",
        "horizontalpodautoscalers/status",
      ],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["storage.k8s.io"],
      resources: [
        "csidrivers",
        "csinodes",
        "storageclasses",
        "volumeattachments",
        "volumeattachments/status",
      ],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["networking.k8s.io"],
      resources: ["networkpolicies"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["scheduling.k8s.io"],
      resources: ["priorityclasses"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["node.k8s.io"],
      resources: ["runtimeclasses"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["extensions"],
      resources: ["ingresses", "ingresses/status"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["events.k8s.io"],
      resources: ["events"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["apiextensions.k8s.io"],
      resources: [
        "customresourcedefinitions",
        "customresourcedefinitions/status",
      ],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["apiregistration.k8s.io"],
      resources: ["apiservices", "apiservices/status"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["discovery.k8s.io"],
      resources: ["endpointslices"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["metrics.k8s.io"],
      resources: ["pods", "nodes"],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["policy"],
      resources: [
        "poddisruptionbudgets",
        "poddisruptionbudgets/status",
        "podsecuritypolicies",
      ],
      verbs: ["get", "list", "watch"],
    },
    {
      apiGroups: ["rbac.authorization.k8s.io"],
      resources: [
        "clusterrolebindings",
        "clusterroles",
        "rolebindings",
        "roles",
      ],
      verbs: ["get", "list", "watch"],
    },
  ],
};

const clusterRoleBindingObj = {
  apiVersion: "rbac.authorization.k8s.io/v1",
  kind: "ClusterRoleBinding",
  metadata: {
    name: "developer-read-all",
  },
  subjects: [
    {
      kind: "ServiceAccount",
      name: "developer",
      namespace: "default",
    },
  ],
  roleRef: {
    kind: "ClusterRole",
    name: "read-all-clusterrole",
    apiGroup: "rbac.authorization.k8s.io",
  },
};

const pspObj = {
  apiVersion: "policy/v1beta1",
  kind: "PodSecurityPolicy",
  metadata: {
    name: "example",
  },
  spec: {
    privileged: false,
    seLinux: {
      rule: "RunAsAny",
    },
    supplementalGroups: {
      rule: "RunAsAny",
    },
    runAsUser: {
      rule: "RunAsAny",
    },
    fsGroup: {
      rule: "RunAsAny",
    },
    volumes: ["*"],
  },
};
