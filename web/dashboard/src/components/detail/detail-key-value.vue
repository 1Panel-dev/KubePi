<template>
  <div>
    <el-link @click="showItem=!showItem">
      {{ showItem ? $t("business.common.pack_up") : $t("business.common.expand") }}
    </el-link>
    <div v-if="showItem">
      <div v-for="(value,key,index) in valueObj" v-bind:key="index" class="myTag">
        <el-tag type="info" size="small" v-if="value.length < 100">
          {{ key }} = {{ value }}
        </el-tag>
        <el-tooltip v-if="value.length > 100" :content="value" placement="top">
          <el-tag type="info" size="small" v-if="value.length >= 100">
            {{ key }} = {{ value.substring(0, 100) + "..." }}
          </el-tag>
        </el-tooltip>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "KoDetailKeyValue",
  components: {},
  props: {
    valueObj: Object,
    show: Boolean
  },
  data () {
    return {
      showItem: false
    }
  },
  methods: {},
  created () {
    this.showItem = this.show ? true : this.showItem
  }
}
</script>

<style scoped>

</style>

