<template>
  <div>
    <el-link @click="showItem=!showItem" v-if="showExpand">
      {{ showItem ? $t("business.common.pack_up") : $t("business.common.expand") }}
    </el-link>
    <div v-if="showItem">
      <div v-for="(value,key,index) in valueObj" v-bind:key="index" class="myTag">
        <el-tag type="info" size="small" v-if="value.length < 50">
          {{ key }} = {{ value }}
        </el-tag>
        <el-tooltip v-if="value.length > 50" :content="value" placement="top">
          <el-tag type="info" size="small" v-if="value.length >= 50">
            {{ key }} = {{ value.substring(0, 50) + "..." }}
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
  },
  data () {
    return {
      showItem: true,
      showExpand: false
    }
  },
  methods: {},
  created () {
    if (this.valueObj) {
      this.showExpand = Object.keys(this.valueObj).length > 5
      this.showItem = !this.showExpand
    }
  }
}
</script>

<style scoped>

</style>

