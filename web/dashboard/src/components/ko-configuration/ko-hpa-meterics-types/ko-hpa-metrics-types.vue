<script>
import KoHpaMetricResource from "@/components/ko-configuration/ko-hpa-meterics-types/ko-hpa-metric-resource"
import KoHpaMetricExternal from "@/components/ko-configuration/ko-hpa-meterics-types/ko-hpa-metric-external"
import KoHpaMetricPods from "@/components/ko-configuration/ko-hpa-meterics-types/ko-hpa-metric-pods"
import KoHpaMetricObject from "@/components/ko-configuration/ko-hpa-meterics-types/ko-hpa-metric-object"

export default {
  name: "KoHpaMetricsTypes",
  functional: true,
  render (createElement, { props, children }) {
    function changeType () {
      const type = props.row.type
      switch (type) {
        case "Resource":
          return KoHpaMetricResource
        case "External":
          return KoHpaMetricExternal
        case "Pods":
          return KoHpaMetricPods
        case "Object":
          return KoHpaMetricObject
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
