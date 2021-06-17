<template>
  <el-table-column class-name="fu-table-operations" :align="align" :width="computeWidth" v-bind="$attrs"
                   v-on="$listeners" v-if="showButtons">
    <template #header>
      {{ $attrs.label }}
      <fu-table-column-select type="dialog" :columns="columns" v-if="columns"/>
    </template>
    <template v-slot:default="{row}">

      <fu-table-button v-for="(btn, i) in defaultButtons" :key="i" v-bind="btn" @click="btn.click(row)"
                       :disabled="disableButton(btn, row)"/>
      <fu-table-more-button :buttons="moreButtons" :row="row" v-if="moreButtons.length > 0"/>
    </template>
  </el-table-column>
</template>

<script>


export default {
  name: "KoTableOperations",
  components: {},
  props: {
    columns: Array,
    align: {
      type: String,
      default: "center"
    },
    width: [String, Number],
    minWidth: [String, Number],
    ellipsis: { // 超过几个按钮时显示省略号，如果只超过一个也不显示省略号
      type: Number,
      default: 1
    },
    buttons: {
      type: Array,
      required: true
    },
  },
  computed: {
    showButtons() {
      return this.buttons?.filter(btn => btn.show !== false)
    },
    hasMore() {
      return this.showButtons?.length > this.ellipsis + 1
    },
    defaultButtons() {
      return this.hasMore ? this.showButtons.slice(0, this.ellipsis) : this.showButtons
    },
    moreButtons() {
      return this.hasMore ? this.showButtons.slice(this.ellipsis) : []
    },
    computeWidth() {
      let buttonsWidth = 20 + this.defaultButtons.length * 38 + 38
      if (this.minWidth) {
        buttonsWidth = buttonsWidth < this.minWidth ? this.minWidth : buttonsWidth
      }
      return this.width ? this.width : buttonsWidth
    },
    disableButton() {
      return (btn, row) => {
        return typeof btn.disabled === "function" ? btn.disabled(row) : btn.disabled
      }
    }
  }
}
</script>
