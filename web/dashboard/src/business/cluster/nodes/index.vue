<template>
  <layout-content header="Nodes">
    <complex-table  :data="data" v-loading="loading" @search="search()">
      <template #toolbar>
        <el-input :placeholder="$t('commons.button.search')" suffix-icon="el-icon-search" clearable v-model="searchName"
                  @change="search(true)" @clear="search(true)"></el-input>
      </template>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" fix max-width="50px">
        <template v-slot:default="{row}">
          <el-link @click="onDetail(row)"> {{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="Internal IP" prop="metadata.name" fix max-width="50px">
        <template v-slot:default="{row}">
          <div v-for="(address,index) in row.status.addresses" v-bind:key="index">
            <span v-if="address.type === 'InternalIP'">{{ address.address }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.node.ready')" prop="status.conditions" fix max-width="50px">
        <template v-slot:default="{row}">
          <div v-for="(condition,index) in row.status.conditions" v-bind:key="index">
            <span v-if="condition.type === 'Ready'">{{ condition.status }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.node.role')" prop="metadata.labels" fix max-width="50px">
        <template v-slot:default="{row}">
          <div v-for="(item,name,index) in row.metadata.labels" :key="index" style="display:inline-block">
            <span v-if="item"></span>
            <el-tag v-if="name.indexOf('node-role.kubernetes.io') > -1">{{ name.substring(name.indexOf(".io") + 4) }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {listNodes} from "@/api/nodes"
import KoTableOperations from "@/components/ko-table-operations"

export default {
  name: "NodeList",
  components: { KoTableOperations, ComplexTable, LayoutContent },
  data () {
    return {
      data: [],
      loading: false,
      searchName: "",
      clusterName: "",
      buttons: [
        {
          label: this.$t("commons.button.view_yaml"),
          icon: "el-icon-view",
          click: (row) => {
            this.$router.push({ path: "/nodes/detail/"+row.metadata.name, query: { yamlShow: "true" } })
          }
        },
      ]
    }
  },
  methods: {
    search () {
      this.loading = true
      listNodes(this.clusterName, this.searchName).then(res => {
        this.loading = false
        this.data = res.items
      })
    },
    onDetail (row) {
      this.$router.push({ path: "/nodes/detail/"+row.metadata.name, query: { yamlShow: "false" } })
    }
  },
  created () {
    this.clusterName = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
