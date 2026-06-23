<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.network.external_name')">
      <el-row>
        <el-col :span="12">
          <el-input :placeholder="$t('commons.placeholder.example_database_domain')" v-model="name" @change="transformation"></el-input>
        </el-col>
      </el-row>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
export default {
  name: "KoServiceExternalName",
  components: { KoCard },
  props: {
    externalName: String
  },
  data () {
    return {
      name: ""
    }
  },
  methods: {
    transformation () {
      this.$emit("update:externalName",this.name)
    }
  },
  created () {
    if (this.externalName) {
      this.name = this.externalName
    }
  }
}
</script>

<style scoped>

</style>
