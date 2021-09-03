<template>
  <div>
    <complex-table :data="events" v-loading="loading">
      <el-table-column :label="$t('business.event.reason')" prop="reason" fix max-width="50px" />
      <el-table-column :label="$t('business.event.type')" prop="type" fix max-width="50px">
        <template v-slot:default="{row}">
          <el-tag v-if="row.type ==='Normal'">{{ row.type }}</el-tag>
          <el-tag v-else type="danger">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.event.message')" prop="message" fix min-width="200px" show-overflow-tooltip />
      <el-table-column :label="$t('commons.table.created_time')" min-width="50px">
        <template v-slot:default="{row}">
          <span>{{ row.eventTime | age }}</span>
        </template>
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import { listEventsWithNsSelector } from "@/api/events"

export default {
  name: "KoDetailEvents",
  components: { ComplexTable },
  props: {
    namespace: String,
    cluster: String,
    selector: String,
  },
  watch: {
    selector: {
      handler(newSelector) {
        if(newSelector) {
          this.search()
        }
      },
      immediate: true,
    },
  },
  data() {
    return {
      loading: false,
      events: [],
    }
  },
  methods: {
    search() {
      listEventsWithNsSelector(this.cluster, this.namespace, this.selector).then((res) => {
        for (const e of res.items) {
          e.eventTime = e.firstTimestamp || e.eventTime || e.lastTimestamp
          this.events.unshift(e)
        }
        this.events.sort(function (a, b) {
          return Date.parse((b.firstTimestamp || b.eventTime || b.lastTimestamp).replace(/-/g, "/")) - Date.parse((a.firstTimestamp || b.eventTime || b.lastTimestamp).replace(/-/g, "/"))
        })
      })
    },
  },
}
</script>

<style scoped>
</style>
