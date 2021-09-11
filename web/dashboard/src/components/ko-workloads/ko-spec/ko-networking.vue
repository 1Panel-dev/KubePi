<template>
  <div>
    <el-form label-position="top" :model="form" :disabled="isReadOnly">
      <el-row :gutter="20">
        <el-col :span="24">
          <el-form-item :label="$t('business.workload.dns_policy')" prop="dnsPolicy">
            <ko-form-item itemType="radio" v-model="form.dnsPolicy" :radios="dns_policy_list" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="6">
          <el-form-item label="hostIPC" prop="hostIPC">
            <ko-form-item itemType="radio" v-model="form.hostIPC" :radios="yes_no_list" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="hostNetwork" prop="hostNetwork">
            <ko-form-item itemType="radio" v-model="form.hostNetwork" :radios="yes_no_list" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="hostPID" prop="hostPID">
            <ko-form-item @change="checkIsValid('hostPID')" itemType="radio" v-model="form.hostPID" :radios="yes_no_list" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="shareProcessNamespace" prop="shareProcessNamespace">
            <ko-form-item @change="checkIsValid('shareProcessNamespace')" itemType="radio" v-model="form.shareProcessNamespace" :radios="yes_no_list" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span=12>
          <span>{{$t('business.workload.nameservers')}}</span>
          <table style="width: 100%;margin-top:10px" class="tab-table">
            <tr v-for="(row,index) in form.dnsConfig.nameservers" :key="index">
              <td width="90%">
                <ko-form-item placeholder="e.g. 1.1.1.1" itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleNameserversDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleNameserversAdd">{{$t('business.workload.add')}}{{$t('business.workload.nameservers')}}</el-button>
              </td>
            </tr>
          </table>
        </el-col>
        <el-col :span=12>
          <span>{{$t('business.workload.searches')}}</span>
          <table style="width: 100%;margin-top:10px" class="tab-table">
            <tr v-for="(row,index) in form.dnsConfig.searches" :key="index">
              <td width="90%">
                <ko-form-item placeholder="e.g. mycompany.com" itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleSearchesDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleSearchesAdd">{{$t('business.workload.add')}}{{$t('business.workload.searches')}}</el-button>
              </td>
            </tr>
          </table>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top:20px">
        <el-col :span=12>
          <label>{{$t('business.workload.dns_resolver')}}</label>
          <table style="width: 100%" class="tab-table">
            <tr>
              <th scope="col" width="43%" align="left"><label>{{$t('business.workload.key')}}</label></th>
              <th scope="col" width="43%" align="left"><label>{{$t('business.workload.value')}}</label></th>
              <th align="left"></th>
            </tr>
            <tr v-for="(row,index) in form.dnsConfig.options" :key="index">
              <td>
                <ko-form-item placeholder="e.g. foo" itemType="input" v-model="row.name" />
              </td>
              <td>
                <ko-form-item placeholder="e.g. bar" itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleOptionsDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleOptionsAdd">{{$t('business.workload.add')}} {{$t('business.workload.dns_resolver')}}</el-button>
              </td>
            </tr>
          </table>
        </el-col>
        <el-col :span=12>
          <label>{{$t('business.workload.host_aliases')}}</label>
          <table style="width: 100%" class="tab-table">
            <tr>
              <th scope="col" width="43%" align="left"><label>{{$t('business.workload.ip_address')}}</label></th>
              <th scope="col" width="43%" align="left"><label>{{$t('business.workload.host_aliases')}}</label></th>
              <th align="left"></th>
            </tr>
            <tr v-for="(row, index) in form.hostAliases" v-bind:key="index">
              <td>
                <ko-form-item placeholder="e.g. 1.1.1.1" itemType="input" v-model="row.ip" />
              </td>
              <td>
                <ko-form-item placeholder="e.g. foo.com,bar.com" itemType="input" v-model="row.hostnames" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleAliasDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleAliasAdd">{{$t('business.workload.add')}}{{$t('business.workload.host_aliases')}}</el-button>
              </td>
            </tr>
          </table>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top:20px">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.hostname')" prop="hostname">
            <ko-form-item placeholder="hostname" itemType="input" v-model="form.hostname" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.sub_domain')" prop="subdomain">
            <ko-form-item placeholder="subdomain" itemType="input" v-model="form.subdomain" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoNetworking",
  components: { KoFormItem },
  props: {
    networkingParentObj: Object,
    isReadOnly: Boolean,
  },
  data() {
    return {
      yes_no_list: [
        { label: this.$t("business.workload.yes"), value: true },
        { label: this.$t("business.workload.no"), value: false },
      ],
      dns_policy_list: [
        { label: "Default", value: "Default" },
        { label: "ClusterFirst", value: "ClusterFirst" },
        { label: "None", value: "None" },
        { label: "ClusterFirstWithHostNet", value: "ClusterFirstWithHostNet" },
      ],
      form: {
        hostIPC: null,
        hostNetwork: null,
        hostPID: null,
        shareProcessNamespace: null,

        subdomain: "",
        dnsPolicy: "",
        hostname: "",
        hostAliases: [],
        dnsConfig: {
          nameservers: [],
          searches: [],
          options: [],
        },
      },
    }
  },
  methods: {
    handleSearchesDelete(index) {
      this.form.dnsConfig.searches.splice(index, 1)
    },
    handleNameserversDelete(index) {
      this.form.dnsConfig.nameservers.splice(index, 1)
    },
    handleOptionsDelete(index) {
      this.form.dnsConfig.options.splice(index, 1)
    },
    handleAliasDelete(index) {
      this.form.hostAliases.splice(index, 1)
    },
    handleNameserversAdd() {
      var item = {
        value: "",
      }
      this.form.dnsConfig.nameservers.push(item)
    },
    handleSearchesAdd() {
      var item = {
        value: "",
      }
      this.form.dnsConfig.searches.push(item)
    },
    handleOptionsAdd() {
      var item = {
        name: "",
        value: "",
      }
      this.form.dnsConfig.options.push(item)
    },
    handleAliasAdd() {
      var item = {
        ip: "",
        hostnames: "",
      }
      this.form.hostAliases.push(item)
    },
    checkIsValid(item) {
      if (this.form.hostPID && this.form.shareProcessNamespace) {
        this.$notify({ title: this.$t("commons.message_box.prompt"), message: "shareProcessNamespace 和 hostPID 不能同时为 true", type: "warning" })
        if (item === "hostPID") {
          this.form.hostPID = null
        } else {
          this.form.shareProcessNamespace = null
        }
      }
    },

    transformation(parentFrom) {
      parentFrom.subdomain = this.form.subdomain || undefined
      parentFrom.dnsPolicy = this.form.dnsPolicy || undefined
      parentFrom.hostname = this.form.hostname || undefined
      parentFrom.hostIPC = this.form.hostIPC === null ? undefined : this.form.hostIPC
      parentFrom.hostNetwork = this.form.hostNetwork === null ? undefined : this.form.hostNetwork
      parentFrom.hostPID = this.form.hostPID === null ? undefined : this.form.hostPID
      parentFrom.shareProcessNamespace = this.form.shareProcessNamespace === null ? undefined : this.form.shareProcessNamespace

      let aliases = []
      for (const item of this.form.hostAliases) {
        aliases.push({
          ip: item.ip || undefined,
          hostnames: item.hostnames ? item.hostnames.split(",") : undefined,
        })
      }
      parentFrom.hostAliases = aliases.length !== 0 ? aliases : undefined

      let nameservers = []
      for (const item of this.form.dnsConfig.nameservers) {
        nameservers.push(item.value)
      }
      let searches = []
      for (const item of this.form.dnsConfig.searches) {
        searches.push(item.value)
      }
      parentFrom.dnsConfig = parentFrom.dnsConfig || {
        nameservers: nameservers.length !== 0 ? nameservers : undefined,
        searches: searches.length !== 0 ? searches : undefined,
        options: this.form.dnsConfig.options.length !== 0 ? this.form.dnsConfig.options : undefined,
      }
    },
  },
  mounted() {
    if (this.networkingParentObj) {
      if (this.networkingParentObj.subdomain) {
        this.form.subdomain = this.networkingParentObj.subdomain
      }
      if (this.networkingParentObj.dnsPolicy) {
        this.form.dnsPolicy = this.networkingParentObj.dnsPolicy
      }
      if (this.networkingParentObj.hostname) {
        this.hostname = this.networkingParentObj.hostname
      }

      if (this.networkingParentObj.hostIPC !== undefined) {
        this.form.hostIPC = this.networkingParentObj.hostIPC
      }
      if (this.networkingParentObj.hostNetwork !== undefined) {
        this.form.hostNetwork = this.networkingParentObj.hostNetwork
      }
      if (this.networkingParentObj.hostPID !== undefined) {
        this.form.hostPID = this.networkingParentObj.hostPID
      }
      if (this.networkingParentObj.shareProcessNamespace !== undefined) {
        this.form.shareProcessNamespace = this.networkingParentObj.hostPID
      }

      if (this.networkingParentObj.hostAliases) {
        this.form.hostAliases = []
        for (const item of this.networkingParentObj.hostAliases) {
          this.form.hostAliases.push({
            ip: item.ip,
            hostnames: item.hostnames.join(","),
          })
        }
      }
      if (this.networkingParentObj.dnsConfig) {
        if (this.networkingParentObj.dnsConfig.nameservers) {
          this.form.dnsConfig.nameservers = []
          for (const item of this.networkingParentObj.dnsConfig.nameservers) {
            this.form.dnsConfig.nameservers.push({ value: item })
          }
        }
        if (this.networkingParentObj.dnsConfig.searches) {
          this.form.dnsConfig.searches = []
          for (const item of this.networkingParentObj.dnsConfig.searches) {
            this.form.dnsConfig.searches.push({ value: item })
          }
        }
        if (this.networkingParentObj.dnsConfig.options) {
          this.form.dnsConfig.options = this.networkingParentObj.dnsConfig.options
        }
      }
    }
  },
}
</script>