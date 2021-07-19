<template>
  <div style="margin-top: 20px">
    <ko-card title="Metrics">
      <el-form label-position="top" ref="form" :model="form">
        <el-card v-for="(row,index) in form.metrics" v-bind:key="index" style="background-color: #292a2e">
          <div>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Source">
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
      <el-button @click="addMetrics"></el-button>
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
      external: {
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
      },
      pods: {
        type: "Pods",
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
      },
      resource: {
        type: "Resource",
        resource: {
          name: "cpu",
          target: {
            type: "Utilization",
            averageUtilization: 80
          }
        }
      },
      refresh: false
    }
  },
  methods: {
    change (row, index) {
      if (row.type === "Resource") {
        row = this.resource
      }
      if (row.type === "External") {
        row = this.external
      }
      if (row.type === "Pods") {
        row = this.pods
      }
      this.form.metrics[index] = row
      this.$emit("update:metricsObj", this.form.metrics)
    },
    addMetrics() {
      this.form.metrics.push(this.resource)
    }
  },
  mounted () {
    if (this.metricsObj) {
      this.form.metrics = this.metricsObj
    }else {
      this.form.metrics.push(this.resource)
      this.$emit("update:metricsObj", this.form.metrics)
    }
  }
}
</script>

<style scoped>

</style>
