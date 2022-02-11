<template>
  <div>
    <complex-table :data="rules">
      <el-table-column label="API Groups">
        <template v-slot:default="{row}">
          <div v-for="(apiGroup,index) in row.apiGroups" :key="index" class="myTag">
            <el-tag v-if="apiGroup!==''" size="mini">{{ apiGroup }}</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="URL">
        <template v-slot:default="{row}">
          <span v-for="(url,index) in row.nonResourceURLs" :key="index">{{ url }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Resources" min-width="300px">
        <template v-slot:default="{row}">
          <el-tag v-for="(resource,index) in row.resources" type="success" :key="index">{{ resource }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Create" width="80px">
        <template v-slot:default="{row}">
          <i v-if="checkVerb(row.verbs,'create')" class="el-icon-check"></i>
          <span v-else> - </span>
        </template>
      </el-table-column>
      <el-table-column label="Delete" width="80px" >
        <template v-slot:default="{row}">
          <i v-if="checkVerb(row.verbs,'delete')" class="el-icon-check"></i>
          <span v-else> - </span>
        </template>
      </el-table-column>
      <el-table-column label="Get" width="80px">
        <template v-slot:default="{row}">
          <i v-if="checkVerb(row.verbs,'get')" class="el-icon-check"></i>
          <span v-else> - </span>
        </template>
      </el-table-column>
      <el-table-column label="List" width="80px">
        <template v-slot:default="{row}">
          <i v-if="checkVerb(row.verbs,'list')" class="el-icon-check"></i>
          <span v-else> - </span>
        </template>
      </el-table-column>
      <el-table-column label="Patch" width="80px">
        <template v-slot:default="{row}">
          <i v-if="checkVerb(row.verbs,'patch')" class="el-icon-check"></i>
          <span v-else> - </span>
        </template>
      </el-table-column>
      <el-table-column label="Update" width="80px">
        <template v-slot:default="{row}">
          <i v-if="checkVerb(row.verbs,'update')" class="el-icon-check"></i>
          <span v-else> - </span>
        </template>
      </el-table-column>
      <el-table-column label="Watch" width="80px">
        <template v-slot:default="{row}">
          <i v-if="checkVerb(row.verbs,'watch')" class="el-icon-check"></i>
          <span v-else> - </span>
        </template>
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"

export default {
  name: "KoDetailRoles",
  components: { ComplexTable },
  props: {
    rules: Array
  },
  data () {
    return {
      data: []
    }
  },
  methods: {
    checkVerb (verbs, ope) {
      return verbs.indexOf("*") > -1 ? verbs.indexOf("*") > -1 : verbs.indexOf(ope) > -1
    }
  },
  created () {

  }
}
</script>

<style scoped>

</style>
