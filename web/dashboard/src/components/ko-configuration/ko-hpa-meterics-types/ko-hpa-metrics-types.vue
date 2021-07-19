<script>
import KoHpaMetricResource from "@/components/ko-configuration/ko-hpa-meterics-types/ko-hpa-metric-resource"
import KoHpaMetricExternal from "@/components/ko-configuration/ko-hpa-meterics-types/ko-hpa-metric-external"

export default {
  name: "KoHpaMetricsTypes",
  functional: true,
  render (createElement, { props, data, children }) {
    function changeType () {
      const type = props.row.type
      switch (type) {
        case "Resource":
          return KoHpaMetricResource
        case "External":
          return KoHpaMetricExternal
        case "Pods":
          return KoHpaMetricExternal
      }
    }

    return createElement(
      changeType(),
      {
        props: {
          row: props.row
        }
      },
      children
    )
  },
  props: {
    type: String,
    row: Object
  }
}
</script>

<style scoped>

</style>
