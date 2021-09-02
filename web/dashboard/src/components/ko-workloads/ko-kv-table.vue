<template>
  <div>
    <el-form :disabled="isReadOnly">
      <el-form-item :label="$t('business.workload.annotations')">
        <table style="width: 100%" class="tab-table">
          <tr v-for="(row,index) in annotations" :key="index">
            <td width="48%">
              <ko-form-item :disabled="disableOperator(row)" placeholder="e.g. foo" clearable itemType="input" v-model="row.key" />
            </td>
            <td width="48%">
              <ko-form-item :disabled="disableOperator(row)" clearable itemType="input" v-model="row.value" />
            </td>
            <td width="4%" v-if="isCreate">
              <el-button type="text" :disabled="disableOperator(row)" style="font-size: 10px" @click="handleDelete(annotations, index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button v-if="isCreate" @click="handleAdd(annotations)">{{ $t("commons.button.add") }}</el-button>
              <el-button v-else @click="handleEdit('annotations')">{{ $t("commons.button.edit") }}</el-button>
            </td>
          </tr>
        </table>
      </el-form-item>
      <el-form-item :label="$t('business.workload.label')" required>
        <table style="width: 100%" class="tab-table">
          <tr v-for="(row,index) in labels" :key="index">
            <td width="48%">
              <ko-form-item :disabled="disableOperator(row)" placeholder="e.g. foo" clearable itemType="input" v-model="row.key" />
            </td>
            <td width="48%">
              <ko-form-item :disabled="disableOperator(row)" clearable itemType="input" v-model="row.value" />
            </td>
            <td width="4%" v-if="isCreate">
              <el-button type="text" :disabled="disableOperator(row)" style="font-size: 10px" @click="handleDelete(labels, index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button v-if="isCreate" @click="handleAdd(labels)">{{ $t("commons.button.add") }}</el-button>
              <el-button v-else @click="handleEdit('labels')">{{ $t("commons.button.edit") }}</el-button>
            </td>
          </tr>
        </table>
      </el-form-item>
    </el-form>

    <el-dialog :title="$t('commons.button.edit')" width="70%" :close-on-click-modal="false" :visible.sync="dialogEditVisible">
      <el-form label-position="top" style="margin-left: 30px" ref="form" :model="form">
        <div v-if="form.isOptionShow">
          <el-form-item :label="$t('business.workload.edit_option')" prop="optionMode">
            <ko-form-item @change="changeMode" itemType="radio" radioLayout='vertical' v-model="form.optionMode" :radios="options" />
            <div v-if="form.optionMode === 'batch'"><span>{{$t('business.workload.notice')}} spec.selector.matchLabels {{$t('business.workload.recreate_tips')}}</span></div>
          </el-form-item>
        </div>
        <table style="width: 100%;" class="tab-table">
          <tr v-for="(row,index) in form.data" :key="index">
            <td width="49%">
              <ko-form-item :disabled="disableDialogRowDelete(row)" placeholder="e.g. foo" clearable itemType="input" v-model="row.key" />
            </td>
            <td width="49%">
              <ko-form-item :disabled="disableDialogRowDelete(row)" clearable itemType="input" v-model="row.value" />
            </td>
            <td width="4%">
              <el-button type="text" :disabled="disableDialogRowDelete(row)" style="font-size: 10px" @click="handleDelete(form.data, index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button @click="handleAdd(form.data)">{{ $t("commons.button.add") }}</el-button>
            </td>
          </tr>
        </table>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="dialogEditVisible = false">{{ $t("commons.button.cancel") }}</el-button>
        <el-button size="small" @click="onSubmit">{{ $t("commons.button.confirm") }}</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"
import { deepClone } from "@/utils/deepClone"

export default {
  name: "KoKvTable",
  components: { KoFormItem },
  props: {
    metadataObj: Object,
    selectorObj: Object,
    resourceName: String,
    isReadOnly: Boolean,
    isCreate: Boolean,
  },
  watch: {
    resourceName: {
      handler(newVal) {
        if (newVal) {
          for (const item of this.labels) {
            if (item.key === "k8s.kubepi.cn/name") {
              item.value = newVal
              break
            }
          }
        }
      },
      immediate: true,
    },
  },
  data() {
    return {
      labels: [],
      annotations: [],
      selectors: [],

      dialogEditVisible: false,
      form: {
        isOptionShow: true,
        optionMode: "single",
        data: [],
        type: "",
      },
      options: [
        { label: this.$t("business.workload.edit") + " metadata.labels", value: "single" },
        { label: this.$t("business.workload.edit") + " metadata.labels / spec.selector.matchLabels / spec.template.metadata.labels", value: "batch" },
      ],
    }
  },
  methods: {
    handleDelete(data, index) {
      data.splice(index, 1)
    },
    handleAdd(data) {
      const item = {
        key: "",
        value: "",
      }
      data.push(item)
    },
    handleEdit(type) {
      switch (type) {
        case "annotations": {
          this.form.data = deepClone(this.annotations)
          this.form.isOptionShow = false
          break
        }
        case "labels": {
          this.form.data = deepClone(this.labels)
          this.form.optionMode = "single"
          this.form.isOptionShow = true
          break
        }
        case "selectors": {
          this.form.data = deepClone(this.selectors)
          this.form.optionMode = "batch"
          this.form.isOptionShow = true
          break
        }
      }
      this.form.type = type
      this.dialogEditVisible = true
    },
    changeMode(val) {
      this.form.data = val === "single" ? deepClone(this.labels) : deepClone(this.selectors)
      this.form.type = val === "single" ? "labels" : "selectors"
      this.$emit("changeCreateMode", true)
    },
    onSubmit() {
      if (this.form.type === "annotations") {
        this.annotations = this.form.data
      } else {
        if (this.form.optionMode === "batch") {
          this.labels = this.form.data
          this.selectors = this.form.data
        } else if (this.form.optionMode === "single") {
          this.labels = this.form.data
        }
      }
      this.dialogEditVisible = false
    },
    disableOperator(row) {
      return !this.isCreate || row.key === "k8s.kubepi.cn/name"
    },
    disableDialogRowDelete(row) {
      return row.key === "k8s.kubepi.cn/name"
    },
    transformation(parentForm, podMetadata) {
      if (this.isCreate) {
        let labels = this.parseArryToObj(this.labels)
        parentForm.spec.selector = { matchLabels: labels }
        podMetadata.labels = labels
        parentForm.metadata.labels = labels
        parentForm.metadata.annotations = this.parseArryToObj(this.annotations)
        return
      }
      let selectors = this.parseArryToObj(this.selectors)
      parentForm.spec.selector = { matchLabels: selectors }
      podMetadata.labels = selectors
      parentForm.metadata.labels = selectors
      parentForm.metadata.annotations = this.parseArryToObj(this.annotations)
      parentForm.metadata.labels = this.parseArryToObj(this.labels)
    },
    parseObjToArry(dataObj) {
      let data = []
      if (dataObj) {
        for (const key in dataObj) {
          if (Object.prototype.hasOwnProperty.call(dataObj, key)) {
            data.push({
              key: key,
              value: dataObj[key],
            })
          }
        }
      }
      return data
    },
    parseArryToObj(data) {
      let obj = {}
      for (let i = 0; i < data.length; i++) {
        if (data[i].key !== "") {
          obj[data[i].key] = data[i].value
        }
      }
      return obj
    },
  },
  mounted() {
    if (this.metadataObj) {
      this.labels = this.parseObjToArry(this.metadataObj.labels)
      this.annotations = this.parseObjToArry(this.metadataObj.annotations)
    }
    if (this.selectorObj) {
      this.selectors = this.parseObjToArry(this.selectorObj.matchLabels)
    }
  },
}
</script>

<style scoped>
</style>
