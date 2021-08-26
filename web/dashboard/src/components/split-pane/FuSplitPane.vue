<template>
  <div :style="{ cursor, userSelect}" class="fu-split-pane" ref="outerWrapper">
    <div :class="[`is-${direction}`,'fu-split-pane__left']" :style="{ [attr]: `${valueL}px`, 'padding-right': padding }">
      <slot :name="isHorizontal ? 'left' : 'top'"></slot>
    </div>
    <div :class="resizerClasses" :style="{ [resizerAttr]: `${valueL}px`, ...resizerStyle}" @mousedown="onMouseDown">
      <div class="icon" v-if="resizerType==='resizer'">
        <slot name="resizer">
          <i class="el-icon-more"></i>
        </slot>
      </div>
    </div>
    <div :class="[`is-${direction}`,'fu-split-pane__right']" :style="{ [attr]: valueR, 'padding-left':padding}">
      <slot :name="isHorizontal ? 'right' : 'bottom'"></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "FuSplitPane",
  props: {
    min: {
      type: [Number, String],
      default: "10px",
    },
    left: [Number, String],
    right: [Number, String],
    top: [Number, String],
    bottom: [Number, String],
    direction: {
      validator: (val) => ["vertical", "horizontal"].includes(val),
      default: "horizontal",
    },
    localKey: String,
    resizable: {
      type: Boolean,
      default: true,
    },
    resizerType: {
      validator: (val) => ["resizer", "line"].includes(val),
      default: "resizer",
    },
    resizerClass: String,
    resizerStyle: Object,
  },
  watch: {
    bottom(newValue) {
      this.valueL = this.getMin(this.outerWrapperSize - this.percentToValue(newValue))
    },
  },
  computed: {
    isHorizontal() {
      return this.direction === "horizontal"
    },
    userSelect() {
      return this.active ? "none" : ""
    },
    cursor() {
      return this.active && this.resizable ? (this.isHorizontal ? "col-resize" : "row-resize") : ""
    },
    outerWrapperSize() {
      return this.screenHeight
    },
    offsetSize() {
      return this.isHorizontal ? "offsetWidth" : "offsetHeight"
    },
    defaultValue() {
      if (this.isHorizontal) {
        return this.left ? this.getMin(this.percentToValue(this.left)) : (this.right && this.getMin(this.outerWrapperSize - this.percentToValue(this.right))) || this.outerWrapperSize / 2
      } else {
        return this.top ? this.getMin(this.percentToValue(this.top)) : (this.bottom && this.getMin(this.outerWrapperSize - this.percentToValue(this.bottom))) || this.outerWrapperSize / 2
      }
    },
    valueR() {
      return `calc(100% - ${this.valueL}px)`
    },
    attr() {
      return this.isHorizontal ? "width" : "height"
    },
    resizerAttr() {
      return this.isHorizontal ? "left" : "top"
    },
    saveKey({ localKey }) {
      return "Fu-SP-" + localKey
    },
    resizerClasses() {
      const classes = [`fu-split-pane__${this.resizerType}`, `is-${this.direction}`, this.resizable && "is-resizable", this.resizerClass]
      return classes
    },
    padding() {
      return this.resizerType === "resizer" && "3px"
    },
  },
  data() {
    return {
      active: false,
      valueL: 0,
      oldValue: 0,
      initOffset: 0,
      screenHeight: null,
    }
  },
  mounted() {
    this.screenHeight = document.body.clientHeight - 55
    window.addEventListener("resize", this.listenResize)
    this.readValue()
  },
  methods: {
    listenResize() {
      this.screenHeight = document.body.clientHeight - 55
      this.valueL = this.getMin(this.screenHeight - this.percentToValue(this.bottom))
    },
    onMouseDown(e) {
      this.initOffset = this.isHorizontal ? e.pageX : e.pageY
      this.oldValue = this.valueL
      this.active = true
      document.addEventListener("mousemove", this.onMouseMove)
      document.addEventListener("mouseup", this.onMouseUp)
    },
    onMouseUp() {
      this.active = false
      document.removeEventListener("mousemove", this.onMouseMove)
      document.removeEventListener("mouseup", this.onMouseUp)
    },
    onMouseMove(e) {
      if (!this.resizable) return
      if (this.active) {
        const currentPage = this.isHorizontal ? e.pageX : e.pageY
        const offset = currentPage - this.initOffset
        const value = this.oldValue + offset
        if (value > this.percentToValue(this.min) && value < this.outerWrapperSize - this.percentToValue(this.min)) {
          this.valueL = value
          this.$store.commit("terminal/CHANGE_TERMINAL_HEIGHT", this.screenHeight - this.valueL)
          this.writeValue()
        }
      }
    },
    // 百分比换算成像素
    percentToValue(val) {
      const size = this.screenHeight
      if (typeof val === "string" && val.includes("%")) {
        return (parseInt(val) / 100) * size
      } else {
        return parseInt(val)
      }
    },
    // 是否取最小值
    getMin(val) {
      return val < this.percentToValue(this.min) ? this.percentToValue(this.min) : val
    },
    // localStorage储存数值
    writeValue() {
      if (this.localKey) {
        localStorage.setItem(this.saveKey, this.valueR)
      }
    },
    readValue() {
      if (this.localKey) {
        const local = localStorage.getItem(this.saveKey)

        this.valueL = this.getMin(this.outerWrapperSize - this.percentToValue(local)) || this.defaultValue
      } else {
        this.valueL = this.defaultValue
      }
    },
  },
}
</script>

