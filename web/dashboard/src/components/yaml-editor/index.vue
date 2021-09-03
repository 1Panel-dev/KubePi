<template>
  <div>
    <div class="yaml-editor">
      <codemirror :merge="diff" ref="editor" v-model="content" :options="options"></codemirror>
    </div>
    <div v-if="!readOnly">
      <el-upload :before-upload="readFile" class="upload" action="">
        <el-button>{{ $t("commons.button.upload") }}</el-button>
      </el-upload>
      <el-button @click="showDiff" v-if="!diff && isEdit" :disabled="content===oldValue" class="upload">
        {{ $t("commons.button.show_diff") }}
      </el-button>
      <el-button @click="diff=false" v-if="diff" class="upload">
        {{ $t("commons.button.continue_edit") }}
      </el-button>
    </div>
  </div>
</template>

<script>
import "codemirror/addon/lint/lint.css"
import "codemirror/lib/codemirror.css"
import "codemirror/theme/ayu-dark.css"
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
import "codemirror/addon/merge/merge"
import "codemirror/addon/merge/merge.css"
import DiffMatchPatch from "diff-match-patch"

window.diff_match_patch = DiffMatchPatch
window.DIFF_DELETE = -1
window.DIFF_INSERT = 1
window.DIFF_EQUAL = 0
window.jsyaml = require("js-yaml")

export default {
  name: "YamlEditor",
  props: {
    value: {
      type: Object,
      default: function () {
        return undefined
      }
    },
    readOnly: {
      type: Boolean,
      default () {
        return false
      }
    },
    isEdit: {
      type: Boolean,
      default () {
        return false
      }
    }
  },
  data () {
    return {
      options: {
        value: "",
        mode: "yaml",
        theme: "ayu-dark",
        lineNumbers: true,
        tabSize: 4,
        foldGutter: true,
        lineWrapping: true,
        gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter", "CodeMirror-lint-markers"],
        lint: true,
        origLeft: null,
        orig: "",
        connect: "align",
        collapseIdentical: true,
        highlightDifference: true,
        // viewportMargin: "Infinity"
      },
      content: "",
      file: File,
      oldValue: "",
      diff: false,
      changed: false
    }
  },
  watch: {
    value: function (newValue) {
      this.initData(newValue)
    }
  },
  mounted () {
    this.initData(this.value)
  },
  methods: {
    getValue () {
      let yaml = require("js-yaml")
      return yaml.load(this.content)
    },
    readFile (file) {
      const reader = new FileReader()
      reader.readAsText(file)
      reader.onerror = e => {
        console.log("error" + e)
      }
      reader.onload = () => {
        this.$refs.editor.codemirror.setValue(reader.result)
      }
      return false
    },
    showDiff () {
      this.diff = !this.diff
      this.options.value = this.content
      this.options.orig = this.oldValue
    },
    changeValue () {
      this.changed = true
    },
    foldCode (key) {
      const cursor = this.$refs.editor.codemirror.getSearchCursor(key)
      if (cursor.findNext()) {
        this.$refs.editor.codemirror.foldCode(cursor.from())
      }
    },
    initData (value) {
      if (value !== undefined) {
        let yaml = require("js-yaml")
        this.oldValue = yaml.dump(value)
        this.$refs.editor.codemirror.setValue(yaml.dump(value))
        //折叠一些无用的key
        const cursor = this.$refs.editor.codemirror.getSearchCursor("managedFields")
        if (cursor.findNext()) {
          this.$refs.editor.codemirror.foldCode(cursor.from())
        }
        this.$refs.editor.codemirror.setOption("readOnly", this.readOnly)
      }
    }
  },
}
</script>

<style lang="scss">
  .yaml-editor {
    height: 100%;
    position: relative;
  }

  .CodeMirror {
    border: 1px solid #eee;
    height: auto;
    min-height: 300px;
  }

  .CodeMirror-merge {
    border: 1px solid #eee;
    height: auto;
    min-height: 300px;
  }

  .CodeMirror-merge-r-chunk {
    background: #3c3c3ce3;
  }

  .CodeMirror-merge-r-chunk-end {
    border-bottom: unset;
  }

  .CodeMirror-merge-r-chunk-start {
    border-top: unset;
  }

  .upload {
    display: inline-block;
    float: left;
    margin-top: 10px;
  }
</style>
