<template>
  <div v-loading="loading">
    <el-dialog title="Metric Server" :visible.sync="dialogMetricVisible" width="50%" @close="onCancel" :close-on-click-modal="false">
      <el-form label-position="top" :model="createForm" ref="createForm" style="margin-left: 20px">
        <el-form-item :label="$t('business.dashboard.metric_server_install_help')">
          <el-tree :data="data" :props="defaultProps"></el-tree>
        </el-form-item>
        <el-form-item :label="$t('business.dashboard.image')" prop="image" :rules="inputRule">
          <el-input v-model="createForm.image" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="onCancel">{{ $t("commons.button.cancel") }}</el-button>
        <el-button :disabled="!hasPermission" size="small" @click="onSubmit">{{ $t("commons.button.confirm") }}</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { checkPermissions } from "@/utils/permission"
import { createService, listServices, updateService } from "@/api/services"
import { createClusterRole, listClusterRoles, updateClusterRole } from "@/api/clusterroles"
import { createClusterRoleBinding, listClusterRoleBindings, updateClusterRoleBinding } from "@/api/clusterrolebindings"
import { createRoleBinding, listRoleBindings, updateRoleBinding } from "@/api/rolebings"
import { createServiceAccount, listServiceAccounts, updateServiceAccount } from "@/api/serviceaccounts"
import { createWorkLoad, listWorkLoads, updateWorkLoad } from "@/api/workloads"
import { createApiServer, listApiServices, updateApiServer } from "@/api/apis"
import Rule from "@/utils/rules"

export default {
  name: "MetricServer",
  props: {
    visible: Boolean,
    clusterName: String,
  },
  data() {
    return {
      loading: false,
      dialogMetricVisible: false,
      hasPermission: true,
      data: [],
      serviceAccountList: [],
      roleBindingList: [],
      clusterRoleList: [],
      clusterRoleBindingList: [],
      serviceList: [],
      deploymentList: [],
      apiServiceList: [],
      defaultProps: {
        children: "children",
        label: "label",
      },
      createForm: {
        image: "",
      },
      inputRule: Rule.RequiredRule,
      form: {
        apiVersion: "v1",
        kind: "List",
        items: [
          {
            apiVersion: "v1",
            kind: "ServiceAccount",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
              },
              name: "metrics-server",
              namespace: "kube-system",
            },
          },
          {
            apiVersion: "rbac.authorization.k8s.io/v1",
            kind: "ClusterRole",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
                "rbac.authorization.k8s.io/aggregate-to-admin": "true",
                "rbac.authorization.k8s.io/aggregate-to-edit": "true",
                "rbac.authorization.k8s.io/aggregate-to-view": "true",
              },
              name: "system:aggregated-metrics-reader",
            },
            rules: [
              {
                apiGroups: ["metrics.k8s.io"],
                resources: ["pods", "nodes"],
                verbs: ["get", "list", "watch"],
              },
            ],
          },
          {
            apiVersion: "rbac.authorization.k8s.io/v1",
            kind: "ClusterRole",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
              },
              name: "system:metrics-server",
            },
            rules: [
              {
                apiGroups: [""],
                resources: ["pods", "nodes", "nodes/stats", "namespaces", "configmaps"],
                verbs: ["get", "list", "watch"],
              },
            ],
          },
          {
            apiVersion: "rbac.authorization.k8s.io/v1",
            kind: "RoleBinding",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
              },
              name: "system:metrics-server",
              namespace: "kube-system",
            },
            roleRef: {
              apiGroups: "rbac.authorization.k8s.io",
              kind: "Role",
              name: "extension-apiserver-authentication-reader",
            },
            subjects: [
              {
                kind: "ServiceAccount",
                name: "metrics-server",
                namespace: "kube-system",
              },
            ],
          },
          {
            apiVersion: "rbac.authorization.k8s.io/v1",
            kind: "ClusterRoleBinding",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
              },
              name: "metrics-server:system:auth-delegator",
              namespace: "kube-system",
            },
            roleRef: {
              apiGroups: "rbac.authorization.k8s.io",
              kind: "ClusterRole",
              name: "system:auth-delegator",
            },
            subjects: [
              {
                kind: "ServiceAccount",
                name: "metrics-server",
                namespace: "kube-system",
              },
            ],
          },
          {
            apiVersion: "rbac.authorization.k8s.io/v1",
            kind: "ClusterRoleBinding",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
              },
              name: "system:metrics-server",
              namespace: "kube-system",
            },
            roleRef: {
              apiGroups: "rbac.authorization.k8s.io",
              kind: "ClusterRole",
              name: "system:metrics-server",
            },
            subjects: [
              {
                kind: "ServiceAccount",
                name: "metrics-server",
                namespace: "kube-system",
              },
            ],
          },
          {
            apiVersion: "v1",
            kind: "Service",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
              },
              name: "metrics-server",
              namespace: "kube-system",
            },
            spec: {
              ports: [{ name: "https", port: 443, protocol: "TCP", targetPort: "https" }],
              selector: {
                "k8s-app": "metrics-server",
              },
            },
          },
          {
            apiVersion: "apps/v1",
            kind: "Deployment",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
              },
              name: "metrics-server",
              namespace: "kube-system",
            },
            spec: {
              selector: {
                matchLabels: {
                  "k8s-app": "metrics-server",
                },
              },
              strategy: {
                rollingUpdate: {
                  maxUnavailable: 0,
                },
              },
              template: {
                metadata: {
                  labels: {
                    "k8s-app": "metrics-server",
                  },
                },
                spec: {
                  containers: [
                    {
                      args: ["--cert-dir=/tmp", "--secure-port=443", "--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname", "--kubelet-use-node-status-port", "--metric-resolution=15s"],
                      image: "kubeoperator/metrics-server:v0.5.0",
                      imagePullPolicy: "IfNotPresent",
                      livenessProbe: {
                        failureThreshold: 3,
                        httpGet: {
                          path: "/livez",
                          port: "https",
                          scheme: "HTTPS",
                        },
                        periodSeconds: 10,
                      },
                      name: "metrics-server",
                      ports: [{ containerPort: 443, name: "https", protocol: "TCP" }],
                      readinessProbe: {
                        failureThreshold: 3,
                        httpGet: {
                          path: "/readyz",
                          port: "https",
                          scheme: "HTTPS",
                        },
                        initialDelaySeconds: 20,
                        periodSeconds: 10,
                      },
                      resources: {
                        requests: {
                          cpu: "100m",
                          memory: "200Mi",
                        },
                      },
                      securityContext: {
                        readOnlyRootFilesystem: true,
                        runAsNonRoot: true,
                        runAsUser: 1000,
                      },
                      volumeMounts: [
                        {
                          mountPath: "/tmp",
                          name: "tmp-dir",
                        },
                      ],
                    },
                  ],
                  nodeSelector: {
                    "kubernetes.io/os": "linux",
                  },
                  priorityClassName: "system-cluster-critical",
                  serviceAccountName: "metrics-server",
                  volumes: [{ emptyDir: {}, name: "tmp-dir" }],
                },
              },
            },
          },
          {
            apiVersion: "apiregistration.k8s.io/v1",
            kind: "APIService",
            metadata: {
              labels: {
                "k8s-app": "metrics-server",
              },
              name: "v1beta1.metrics.k8s.io",
            },
            spec: {
              group: "metrics.k8s.io",
              groupPriorityMinimum: 100,
              insecureSkipTLSVerify: true,
              service: {
                name: "metrics-server",
                namespace: "kube-system",
              },
              version: "v1beta1",
              versionPriority: 100,
            },
          },
        ],
      },
    }
  },
  methods: {
    onSubmit() {
      this.loading = true
      this.$refs["createForm"].validate((valid) => {
        if (valid) {
          let ps = []
          for (const item of this.form.items) {
            let isExist = false
            switch (item.kind.toLowerCase() + "s") {
              case "serviceaccounts":
                for (const s of this.serviceAccountList) {
                  if (s.metadata.name === item.metadata.name) {
                    isExist = true
                    break
                  }
                }
                if (isExist) {
                  ps.push(updateServiceAccount(this.clusterName, item.metadata.namespace, item.metadata.name, item))
                } else {
                  ps.push(createServiceAccount(this.clusterName, item.metadata.namespace, item))
                }
                break
              case "rolebindings":
                for (const s of this.roleBindingList) {
                  if (s.metadata.name === item.metadata.name) {
                    isExist = true
                    break
                  }
                }
                if (isExist) {
                  ps.push(updateRoleBinding(this.clusterName, item.metadata.namespace, item.metadata.name, item))
                } else {
                  ps.push(createRoleBinding(this.clusterName, item.metadata.namespace, item))
                }
                break
              case "clusterroles":
                for (const s of this.clusterRoleList) {
                  if (s.metadata.name === item.metadata.name) {
                    isExist = true
                    break
                  }
                }
                if (isExist) {
                  ps.push(updateClusterRole(this.clusterName, item.metadata.name, item))
                } else {
                  ps.push(createClusterRole(this.clusterName, item))
                }
                break
              case "clusterrolebindings":
                for (const s of this.clusterRoleBindingList) {
                  if (s.metadata.name === item.metadata.name) {
                    isExist = true
                    break
                  }
                }
                if (isExist) {
                  ps.push(updateClusterRoleBinding(this.clusterName, item.metadata.name, item))
                } else {
                  ps.push(createClusterRoleBinding(this.clusterName, item))
                }
                break
              case "services":
                for (const s of this.serviceList) {
                  if (s.metadata.name === item.metadata.name) {
                    isExist = true
                    item.metadata.resourceVersion = s.metadata.resourceVersion
                    break
                  }
                }
                if (isExist) {
                  ps.push(updateService(this.clusterName, item.metadata.namespace, item.metadata.name, item))
                } else {
                  ps.push(createService(this.clusterName, item.metadata.namespace, item))
                }
                break
              case "deployments":
                for (const s of this.deploymentList) {
                  if (s.metadata.name === item.metadata.name) {
                    isExist = true
                    break
                  }
                }
                item.spec.template.spec.containers[0].image = this.createForm.image
                if (isExist) {
                  ps.push(updateWorkLoad(this.clusterName, "deployments", item.metadata.namespace, item.metadata.name, item))
                } else {
                  ps.push(createWorkLoad(this.clusterName, "deployments", item.metadata.namespace, item))
                }
                break
              case "apiservices":
                for (const s of this.apiServiceList) {
                  if (s.metadata.name === item.metadata.name) {
                    isExist = true
                    break
                  }
                }
                if (isExist) {
                  ps.push(updateApiServer(this.clusterName, item.metadata.name, item))
                } else {
                  ps.push(createApiServer(this.clusterName, item))
                }
                break
            }
          }
          Promise.all(ps)
            .then(() => {
              this.$message({
                type: "success",
                message: this.$t("commons.msg.operation_success"),
              })
              this.loading = false
              this.onCancel()
            })
            .catch(() => {
              this.loading = false
            })
        }
      })
    },
    onCancel() {
      this.dialogMetricVisible = false
      this.$emit("changeVisble", this.dialogMetricVisible)
    },
    loadDatas() {
      listServiceAccounts(this.clusterName, "kube-system").then((res) => {
        this.serviceAccountList = res.items
      })
      listRoleBindings(this.clusterName, "kube-system").then((res) => {
        this.roleBindingList = res.items
      })
      listClusterRoles(this.clusterName).then((res) => {
        this.clusterRoleList = res.items
      })
      listClusterRoleBindings(this.clusterName).then((res) => {
        this.clusterRoleBindingList = res.items
      })
      listServices(this.clusterName, "kube-system").then((res) => {
        this.serviceList = res.items
      })
      listWorkLoads(this.clusterName, "deployments", "kube-system").then((res) => {
        this.deploymentList = res.items
      })
      listApiServices(this.clusterName).then((res) => {
        this.apiServiceList = res.items
      })
    },
    yamlInit() {
      for (const item of this.form.items) {
        let isExist = false
        for (const d of this.data) {
          if (d.label === item.kind) {
            isExist = true
            d.count++
            d.children.push({ label: item.metadata.name })
            break
          }
        }
        if (!isExist) {
          let hasPermission = false
          switch (item.kind) {
            case "ServiceAccount":
              hasPermission = checkPermissions({ scope: "namespace", apiGroup: "", resource: "serviceaccounts", verb: "*" })
              break
            case "RoleBinding":
              hasPermission = checkPermissions({ scope: "namespace", apiGroup: "rbac.authorization.k8s.io", resource: "rolebindings", verb: "*" })
              break
            case "ClusterRole":
              hasPermission = checkPermissions({ scope: "cluster", apiGroup: "rbac.authorization.k8s.io", resource: "clusterroles", verb: "*" })
              break
            case "ClusterRoleBinding":
              hasPermission = checkPermissions({ scope: "cluster", apiGroup: "rbac.authorization.k8s.io", resource: "clusterrolebindings", verb: "*" })
              break
            case "Service":
              hasPermission = checkPermissions({ scope: "namespace", apiGroup: "", resource: "services", verb: "*" })
              break
            case "Deployment":
              this.createForm.image = item.spec.template.spec.containers[0].image
              hasPermission = checkPermissions({ scope: "namespace", apiGroup: "apps", resource: "deployments", verb: "*" })
              break
            case "APIService":
              hasPermission = checkPermissions({ scope: "cluster", apiGroup: "apiregistration.k8s.io/v1", resource: "apiservices", verb: "*" })
              break
          }
          this.data.push({ label: item.kind, count: 1, hasPermission: hasPermission, children: [{ label: item.metadata.name }] })
        }
      }
      for (const item of this.data) {
        if (!item.hasPermission) {
          this.hasPermission = false
        }
        item.label = item.label + "  [" + item.count + "]    ---" + (item.hasPermission ? this.$t("business.dashboard.has_permission") : this.$t("business.dashboard.has_not_permission"))
      }
      this.loadDatas()
    },
  },
  created() {
    this.yamlInit()
    this.dialogMetricVisible = true
  },
}
</script>