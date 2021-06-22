<template>
  <div>
    <div class="yaml-editor">
      <codemirror ref="editor" v-model="content" :options="options"></codemirror>
    </div>
  </div>
</template>


<script>
import "codemirror/addon/lint/lint.css"
import "codemirror/lib/codemirror.css"
import "codemirror/theme/erlang-dark.css"
import "codemirror/mode/yaml/yaml"
import "codemirror/addon/lint/lint"
import "codemirror/addon/lint/yaml-lint"
import "codemirror/addon/fold/foldgutter.css"
import "codemirror/addon/fold/foldcode"
import "codemirror/addon/fold/foldgutter"
import "codemirror/addon/fold/brace-fold"
import "codemirror/addon/fold/comment-fold"
import "codemirror/addon/fold/indent-fold"
import "codemirror/addon/search/searchcursor"
import "codemirror/addon/search/search"
import "js-yaml"

window.jsyaml = require("js-yaml")
export default {
  name: "YamlEditor",
  props: {
    value: {
      type: Object,
      default: {}
    },
    readOnly: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      options: {
        value: "",
        mode: "yaml",
        theme: "erlang-dark",
        lineNumbers: true,
        tabSize: 4,
        foldGutter: true,
        lineWrapping: true,
        gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter", "CodeMirror-lint-markers"],
        lint: true
      },
      content: ""
    }
  },
  watch: {},
  mounted () {
    if (this.value !== undefined) {
      let yaml = require("js-yaml")
      this.$refs.editor.codemirror.setValue(yaml.dump(this.value))
      //折叠一些无用的key
      const cursor = this.$refs.editor.codemirror.getSearchCursor("managedFields")
      if (cursor.findNext()) {
        this.$refs.editor.codemirror.foldCode(cursor.from())
      }
      this.$refs.editor.codemirror.setOption("readOnly", this.readOnly)
    }
  },
  methods: {
    getValue () {
      let yaml = require("js-yaml")
      return yaml.load(this.content)
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
        min-height: 600px;
    }

    .yaml-editor >>> .CodeMirror-scroll {
        min-height: 600px;
    }
</style>
