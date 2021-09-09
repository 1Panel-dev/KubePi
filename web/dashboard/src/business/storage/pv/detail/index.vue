<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'PersistentVolumes'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="24">
        <el-col :span="24">
          <el-card class="el-card">
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" align="left">
                  <h3>{{ $t("business.common.basic") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>{{ $t("commons.table.name") }}</td>
                <td colspan="3">{{ item.metadata.name }}</td>
              </tr>
              <tr>
                <td>{{ $t('commons.table.status')}}</td>
                <td colspan="3">
                  <el-button v-if="item.status.phase === 'Bound'" type="success" size="mini" plain round>
                    {{ item.status.phase }}
                  </el-button>
                  <el-button v-else type="" size="mini" plain round>
                    {{ item.status.phase }}
                  </el-button>
                </td>
              </tr>
              <tr>
                <td>{{ this.$t("business.storage.capacity") }}</td>
                <td colspan="3">{{ item.spec.capacity.storage }}</td>
              </tr>
              <tr>
                <td>{{ this.$t("business.storage.accessModes") }}</td>
                <td colspan="3">
                  <div v-for="(value,index) in item.spec.accessModes" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small">
                      {{ value }}
                    </el-tag>
                  </div>
                </td>
              </tr>
              <tr>
                <td>{{ this.$t("business.storage.storageClass") }}</td>
                <td colspan="3">{{ item.spec.storageClassName }}</td>
              </tr>
              <tr>
                <td>{{ this.$t("business.storage.reclaimPolicy") }}</td>
                <td colspan="3">{{ item.spec.persistentVolumeReclaimPolicy }}</td>
              </tr>
              <tr>
                <td>{{ $t("commons.table.created_time") }}</td>
                <td colspan="">{{ item.metadata.creationTimestamp | datetimeFormat }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.common.annotation") }}</td>
                <td colspan="3">
                  <div v-for="(value,key,index) in item.metadata.annotations" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small" v-if="value.length < 100">
                      {{ key }} = {{ value }}
                    </el-tag>
                    <el-tooltip v-if="value.length > 100" :content="value" placement="top">
                      <el-tag type="info" size="small" v-if="value.length >= 100">
                        {{ key }} = {{ value.substring(0, 100) + "..." }}
                      </el-tag>
                    </el-tooltip>
                  </div>
                </td>
              </tr>
            </table>
            <div class="bottom-button">
              <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.view_yaml") }}</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
      <el-row style="margin-top: 20px" class="row-box" :gutter="10">
        <el-col :span="8">
          <el-card v-if="item.spec.claimRef !== undefined">
            <div class="card_title">
              <h3>{{ this.$t("business.storage.claim") }}</h3>
            </div>
            <table style="width: 100%" class="myTable">
              <tr>
                <td>{{ this.$t("business.configuration.type") }}</td>
                <td colspan="3">{{ item.spec.claimRef.kind }}</td>
              </tr>
              <tr>
                <td>{{ this.$t("commons.table.name") }}</td>
                <td colspan="3">{{ item.spec.claimRef.name }}</td>
              </tr>
              <tr>
                <td>{{ this.$t("business.namespace.namespace") }}</td>
                <td colspan="3">{{ item.spec.claimRef.namespace }}</td>
              </tr>
            </table>
          </el-card>
        </el-col>
        <el-col :span="16">
          <el-card v-if="'nfs' in item.spec">
            <div class="card_title">
              <h3>Network File System</h3>
            </div>
            <table style="width: 100%" class="myTable">
              <tr>
                <td>Server</td>
                <td colspan="3">{{ item.spec.nfs.server }}</td>
              </tr>
              <tr>
                <td>Path</td>
                <td colspan="3">{{ item.spec.nfs.path }}</td>
              </tr>
              <tr v-if="item.spec.nfs.readOnly">
                <td>Read Only</td>
                <td colspan="3">{{ item.spec.nfs.readOnly }}</td>
              </tr>
            </table>
          </el-card>
        </el-col>
      </el-row>
    </div>
    <el-row>
      <div v-if="yamlShow">
        <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
        <div class="bottom-button">
          <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
        </div>
      </div>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { isJSON } from "@/utils/data"
import { getPv } from "@/api/pv"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "PersistentVolumeDetail",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
  },
  data() {
    return {
      item: {
        metadata: {},
        spec: {
          capacity: {
            storage: "",
          },
        },
        status: {},
      },
      cluster: "",
      yamlShow: false,
      loading: false,
      yaml: {},
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getPv(this.cluster, this.name).then((res) => {
        this.loading = false
        this.item = res
        if (this.yamlShow) {
          this.yaml = JSON.parse(JSON.stringify(this.item))
        }
      })
    },
    getContent(value) {
      const { Base64 } = require("js-base64")
      const content = Base64.decode(value)
      return JSON.parse(content)
    },
    jsonV(str) {
      const { Base64 } = require("js-base64")
      const content = Base64.decode(str)
      return isJSON(content)
    },
    getValue(value) {
      const { Base64 } = require("js-base64")
      return Base64.decode(value)
    },
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "PersistentVolumeDetail",
        params: { name: this.name },
        query: { yamlShow: newValue },
      })
      this.getDetail()
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
