<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <el-row>
      <el-col :span="20" :offset="2">
        <div>
          <table style="width: 90%" class="myTable">
            <tr>
              <th scope="col" width="30%"></th>
              <th scope="col"></th>
            </tr>
            <tr>
              <td>Name</td>
              <td>{{ item.metadata.name }}</td>
            </tr>
            <tr>
              <td>Labels</td>
              <td>
                <div v-for="(value,key,index) in item.metadata.labels" v-bind:key="index">
                  <el-tag type="info" size="small">
                    {{ key }} = {{ value }}
                  </el-tag>
                </div>
              </td>
            </tr>
            <tr>
              <td>Annotations</td>
              <td>
                <div v-for="(value,key,index) in item.metadata.annotations" v-bind:key="index">
                  <el-tag type="info" size="small">
                    {{ key }} = {{ value.length > 100 ? value.substring(0, 100) + "..." : value }}
                  </el-tag>
                </div>
              </td>
            </tr>
            <tr>
              <td>Finalizers</td>
              <td>
                <div v-for="value in item.metadata.finalizers" v-bind:key="value">
                  <el-tag type="info" size="small">
                    {{ value.length > 100 ? value.substring(0, 100) + "..." : value }}
                  </el-tag>
                </div>
              </td>
            </tr>
            <tr>
              <td>Status</td>
              <td>
                <el-button v-if="item.status.phase ==='Active'" type="success" size="mini" plain round>
                  {{ item.status.phase }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td>Created</td>
              <td>
                {{ item.metadata.creationTimestamp | datetimeFormat }}
              </td>
            </tr>
          </table>
          <div style="float: right;margin-top: 10px">
            <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
          </div>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getNamespace} from "@/api/namespaces"

export default {
  name: "NamespaceDetail",
  components: { LayoutContent },
  props: {
    name: String,
    cluster: String
  },
  data () {
    return {
      item: {
        metadata: {},
        status: {}
      },
      loading: false
    }
  },
  methods: {
    getNamespaceByName () {
      this.loading = true
      getNamespace(this.cluster, this.name).then(res => {
        this.loading = false
        this.item = res
      })
    },
    onCancel(){
      this.$router.push({ name: "Namespaces" })
    }
  },
  created () {
    this.getNamespaceByName()
  },
}
</script>

<style scoped>

</style>