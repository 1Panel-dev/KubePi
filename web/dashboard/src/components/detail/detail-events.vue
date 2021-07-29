<template>
  <div>
    <complex-table :data="events" v-loading="loading">
      <el-table-column :label="$t('business.event.reason')" prop="reason" fix max-width="50px">
      </el-table-column>
      <el-table-column :label="$t('business.event.type')" prop="type" fix max-width="50px">
        <template v-slot:default="{row}">
          <el-tag v-if="row.type ==='Normal'">{{ row.type }}</el-tag>
          <el-tag v-else type="danger">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace" fix max-width="50px">

      </el-table-column>
      <el-table-column :label="$t('business.event.message')" prop="message" fix min-width="200px"
                       show-overflow-tooltip>
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import {listEventsWithNsSelector} from "@/api/events"

export default {
  name: "KoDetailEvents",
  components: { ComplexTable },
  props: {
    namespace: String,
    cluster: String,
    selector: String
  },
  data () {
    return {
      loading: false,
      events: []
    }
  },
  methods: {
    list () {
      listEventsWithNsSelector(this.cluster, this.namespace, this.selector).then(res => {
        this.events = res.items
      })
    }
  },
  created () {
    this.list()
  }
}
</script>

<style scoped>

</style>
