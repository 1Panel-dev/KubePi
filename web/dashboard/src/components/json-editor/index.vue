<template>
  <div class="yaml-editor">
    <h4>{{ key }}</h4>
    <codemirror ref="editor" v-model="content" :options="options"></codemirror>
  </div>
</template>

<script>
import "codemirror/lib/codemirror.css"
import "codemirror/theme/ayu-dark.css"
import "codemirror/mode/javascript/javascript"

export default {
  name: "JsonEditor",
  props: {
    value: {
      type: Object,
    },
    key: {
      type: String,
      default: "",
    }
  },
  data () {
    return {
      options: {
        value: "",
        mode: "application/json",
        theme: "ayu-dark",
        lineNumbers: true,
        tabSize: 2,
        readOnly: true,
        lineWrapping: true,
        gutters: ["CodeMirror-lint-markers"],
      },
      content: "",
      file: File
    }
  },
  watch: {
    value: function (newValue) {
      let content = JSON.stringify(newValue, null, "\t")
      this.$refs.editor.codemirror.setValue(content)
      if (content.indexOf("\\n") > -1) {
        this.$refs.editor.codemirror.setOption("lineSeparator", "\\n")
      }
    }
  },
  mounted () {
    let content = JSON.stringify(this.value, null, "\t")
    this.$refs.editor.codemirror.setValue(content)
    if (content.indexOf("\\n") > -1) {
      this.$refs.editor.codemirror.setOption("lineSeparator", "\\n")
    }
  }
}
</script>

<style scoped>
    .yaml-editor {
        height: 100%;
        position: relative;
    }

    .yaml-editor >>> .CodeMirror {
        height: auto;
        min-height: 300px;
    }
</style>
