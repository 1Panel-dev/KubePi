<template>
  <div>
    <ko-card :title="$t('business.configuration.metrics')">
      <el-form label-position="top" ref="form" :model="form">
        <el-card v-for="(row,index) in form.metrics" v-bind:key="index"
                 style="background-color: #292a2e;margin-top: 10px;position: relative">
          <div style="float: right; padding: 3px 0;position: relative;z-index: 1">
            <el-button  type="text" v-if="form.metrics.length > 1"
                        @click="removeMetrics(index)">{{ $t("commons.button.delete") }}
            </el-button>
          </div>
          <div>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('business.configuration.source')">
                  <el-select v-model="row.type" style="width: 100%" value-key="type" @change="change(row,index)">
                    <el-option label="Resource" value="Resource"></el-option>
                    <el-option label="External" value="External"></el-option>
                    <el-option label="Object" value="Object"></el-option>
                    <el-option label="Pods" value="Pods"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
            <ko-hpa-metrics-types :row="form.metrics[index]"></ko-hpa-metrics-types>
          </div>
        </el-card>
      </el-form>
      <el-button @click="addMetrics" style="margin-top: 10px">{{ $t("commons.button.add") }}</el-button>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
import KoHpaMetricsTypes from "@/components/ko-configuration/ko-hpa-meterics-types/ko-hpa-metrics-types"

export default {
  name: "KoHpaMetrics",
  components: { KoHpaMetricsTypes, KoCard },
  props: {
    metricsObj: Array
  },
  data () {
    return {
      form: {
        metrics: [],
      },
      external: {},
      pods: {},
      resource: {},
      object: {},
    }
  },
  methods: {
    change (row, index) {
      if (row.type === "Resource") {
        this.resource = {
          type: "Resource",
          resource: {
            name: "cpu",
            target: {
              type: "Utilization",
              averageUtilization: 80
            }
          }
        }
        row = this.resource
      }
      if (row.type === "External") {
        this.external = {
          type: "External",
          external: {
            type: "io.k8s.api.autoscaling.v2beta2.externalmetricsource",
            metric: {
              name: "",
              selector: {
                matchExpressions: []
              }
            },
            target: {
              type: "AverageValue",
              averageValue: 80
            }
          }
        }
        row = this.external
      }
      if (row.type === "Pods") {
        this.pods =  {
          type: "Pods",
          pods: {
            type: "io.k8s.api.autoscaling.v2beta2.podsmetricsource",
            metric: {
              name: "",
              selector: {
                matchExpressions: []
              }
            },
            target: {
              type: "AverageValue",
              averageValue: 80
            }
          }
        }
        row = this.pods
      }
      if (row.type === "Object") {
        this.object = {
          type: "Object",
          object: {
            type: "io.k8s.api.autoscaling.v2beta2.objectmetricsource",
            metric: {
              name: "",
              selector: {
                matchExpressions: []
              }
            },
            target: {
              type: "AverageValue",
              averageValue: 80
            },
            describedObject: {
              apiVersion: "",
              kind: "",
              name: ""
            }
          }
        }
        row = this.object
      }
      this.form.metrics[index] = row
      this.$emit("update:metricsObj", this.form.metrics)
    },
    addMetrics () {
      this.addResource()
    },
    removeMetrics (index) {
      this.form.metrics.splice(index, 1)
    },
    addResource() {
      this.resource = {
        type: "Resource",
        resource: {
          name: "cpu",
          target: {
            type: "Utilization",
            averageUtilization: 80
          }
        }
      }
      this.form.metrics.push(this.resource)
    },
  },
  mounted () {
    if (this.metricsObj) {
      this.form.metrics = this.metricsObj
    } else {
      this.addResource()
      this.$emit("update:metricsObj", this.form.metrics)
    }
  }
}
</script>

<style scoped>

</style>
