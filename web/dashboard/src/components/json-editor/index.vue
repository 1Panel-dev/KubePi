<template>
  <div class="yaml-editor">
    <codemirror ref="editor" v-model="content" :options="options"></codemirror>
  </div>
</template>

<script>
import "codemirror/lib/codemirror.css"
import "codemirror/theme/material-darker.css"
import "codemirror/mode/javascript/javascript"
import "js-yaml"

export default {
  name: "JsonEditor",
  props: {
    value: {
      type: Object,
      default: function () {
        return {}
      }
    },
  },
  data () {
    return {
      options: {
        value: "",
        mode: {
          name: "javascript",
          json: true,
          statementIndent: 2
        },
        theme: "material-darker",
        lineNumbers: true,
        tabSize: 2,
        foldGutter: true,
        lineWrapping: true,
        lint: true,
        readOnly: true
      },
      content: "",
      file: File
    }
  },
  watch: {
    value: function () {
      let content = JSON.stringify(this.value, null, "\t")
      this.$refs.editor.codemirror.setValue(content)
      this.$refs.editor.codemirror.setOption("lineSeparator", '\\n')
    }
  },
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
        /*background-color: #2c2c2c;*/
    }
</style>