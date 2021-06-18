<template>
  <layout-content header="Nodes">
    <complex-table :pagination-config="page" :data="data" v-loading="loading" @search="search()">
      <template #toolbar>
        <el-input :placeholder="$t('commons.button.search')" suffix-icon="el-icon-search" clearable v-model="searchName"
                  @change="search(true)" @clear="search(true)"></el-input>
      </template>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" fix max-width="50px">
        <template v-slot:default="{row}">
          <el-link @click="onDetail(row)"> {{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="Internal IP"  prop="metadata.name" fix max-width="50px">
        <template v-slot:default="{row}">
          <div v-for="(address,index) in row.status.addresses" v-bind:key="index">
            <span v-if="address.type === 'InternalIP'">{{address.address}}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.node.ready')" prop="status.conditions" fix max-width="50px">
        <template v-slot:default="{row}">
          <div v-for="(condition,index) in row.status.conditions" v-bind:key="index">
            <span v-if="condition.type === 'Ready'">{{condition.status}}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.node.role')" prop="metadata.labels" fix max-width="50px">
        <template v-slot:default="{row}">
          <div v-for="(item,name,index) in row.metadata.labels" :key="index">
            <span v-if="item"></span>
            <el-tag v-if="name.indexOf('node-role.kubernetes.io') > -1">{{name.substring(name.indexOf('.io')+4)}}</el-tag>
          </div>
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {listNodes} from "@/api/nodes"

export default {
  name: "NodeList",
  components: { ComplexTable, LayoutContent },
  data () {
    return {
      data: [],
      page: {
        pageSize: 10,
        nextToken: "",
      },
      loading: false,
      searchName: "",
      clusterName: "test1",
    }
  },
  methods: {
    search (init) {
      this.loading = true
      if (init) {
        this.page = {
          pageSize: this.page.pageSize,
          nextToken: "",
        }
      }
      listNodes(this.clusterName, this.page.pageSize, this.page.nextToken, this.searchName).then(res => {
        this.loading = false
        this.data = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
      })
    },
    onDetail(row) {
      this.$router.push({ name: "NodeDetail",params:{name:row.metadata.name,cluster:this.clusterName} })
    }
  },
  created () {
    this.search()
  }
}
</script>

<style scoped>

</style>