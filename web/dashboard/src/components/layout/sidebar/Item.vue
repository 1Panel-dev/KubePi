<script>
import Namespaced from "@/components/layout/sidebar/namespaced"
import NamespacedTip from "@/components/layout/sidebar/namespaced"
export default {
  name: "MenuItem",
  components: { NamespacedTip, Namespaced },
  functional: true,
  props: {
    icon: {
      type: String,
      default: ""
    },
    title: {
      type: String,
      default: ""
    },
    namespaced: {
      type: Boolean,
      default: false
    }
  },
  render (h, context) {
    const { icon, title, namespaced } = context.props
    const vnodes = []

    if (icon) {
      vnodes.push(<i class={[icon, "sub-el-icon"]}/>)
    }

    if (title &&namespaced) {
      vnodes.push(<span slot="title">{(title)}<namespaced-tip></namespaced-tip></span> )
    }else {
      vnodes.push(<span slot="title">{(title)}</span>)
    }

    return vnodes
  }
}
</script>

<style scoped>
    .sub-el-icon {
        font-size: 20px;
        width: 20px;
        height: 20px;
    }
</style>
