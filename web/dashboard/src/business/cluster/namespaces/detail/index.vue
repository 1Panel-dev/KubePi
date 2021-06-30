<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <el-row>
      <el-card v-if="!yamlShow">
        <el-col :span="24">
            <table class="myTable">
              <tr>
                <th scope="col" width="30%"  align="left">
                  <h3>{{ $t("business.common.basic") }}</h3>
                </th>
                <th scope="col">
                </th>
              </tr>
              <tr>
                <td>{{ $t("commons.table.name") }}</td>
                <td>{{ item.metadata.name }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.common.label") }}</td>
                <td>
                  <div v-for="(value,key,index) in item.metadata.labels" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small">
                      {{ key }} = {{ value }}
                    </el-tag>
                  </div>
                </td>
              </tr>
              <tr>
                <td>{{ $t("business.common.annotation") }}</td>
                <td colspan="4">
                  <div v-for="(value,key,index) in item.metadata.annotations" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small" v-if="value.length < 100">
                      {{ key }} = {{value}}
                    </el-tag>
                    <el-tooltip  v-if="value.length > 100" :content="value" placement="top">
                      <el-tag type="info" size="small" v-if="value.length >= 100" >
                        {{ key }} = {{value.substring(0, 100) + "..." }}
                      </el-tag>
                    </el-tooltip>
                  </div>
                </td>
              </tr>
              <tr>
                <td>Finalizers</td>
                <td>
                  <div v-for="value in item.metadata.finalizers" v-bind:key="value" class="myTag">
                    <el-tag type="info" size="small">
                      {{ value.length > 100 ? value.substring(0, 100) + "..." : value }}
                    </el-tag>
                  </div>
                </td>
              </tr>
              <tr>
                <td>{{ $t("commons.table.status") }}</td>
                <td>
                  <el-button v-if="item.status.phase ==='Active'" type="success" size="mini" plain round>
                    {{ item.status.phase }}
                  </el-button>
                </td>
              </tr>
              <tr>
                <td>{{ $t("commons.table.created_time") }}</td>
                <td>
                  {{ item.metadata.creationTimestamp | datetimeFormat }}
                </td>
              </tr>
            </table>
          <div class="bottom-button">
            <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.view_yaml") }}</el-button>
          </div>
        </el-col>
      </el-card>
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
import {getNamespace} from "@/api/namespaces"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "NamespaceDetail",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String
  },
  data () {
    return {
      item: {
        metadata: {},
        status: {}
      },
      loading: false,
      yamlShow: false,
      yaml: {}
    }
  },
  methods: {
    getNamespaceByName () {
      this.loading = true
      getNamespace(this.cluster, this.name).then(res => {
        this.loading = false
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
      })
    },
  },
  watch: {
    yamlShow:function (newValue) {
      this.$router.push({
        path: "/namespaces/detail/"+this.name,
        query: { yamlShow: newValue }
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === 'true'
    this.getNamespaceByName()
  },
}
</script>

<style scoped>

</style>