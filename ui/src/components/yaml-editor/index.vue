<template>
  <div>
    <div class="yaml-editor">
      <codemirror :value="value" :options="options"></codemirror>
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
import "js-yaml"
export default {
  name: "YamlEditor",
  props: ["value"],
  data () {
    return {
      options: {
        value:'',
        mode:"text/yaml",
        theme: "erlang-dark",
        lineNumbers: true,
        tabSize:4
      }
    }
  },
  watch: {

  },
  mounted () {
      if(this.value !== "" ) {
        let yaml = require("js-yaml");
        this.value = yaml.dump(this.value)
      }
  },
  methods: {

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

    .yaml-editor >>> .CodeMirror-scroll {
        min-height: 300px;
    }

    .yaml-editor >>> .cm-s-rubyblue span.cm-string {
        color: #141418;
    }
</style>
