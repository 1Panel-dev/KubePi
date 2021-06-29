<template>
  <div style="margin-top: 20px">
    <ko-card title="Networking">
      <el-row :gutter="20">
        <el-col :span="12">
          <ko-form-item labelName="Network Mode" itemType="select" v-model="form.hostNetwork" :selections="network_mode_list" />
        </el-col>
        <el-col :span="12">
          <ko-form-item labelName="DNS Policy" itemType="select" v-model="form.dnsPolicy" :selections="dns_policy_list" />
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <ko-form-item labelName="Hostname" placeholder="e.g. web" itemType="input" v-model="form.hostname" />
        </el-col>
        <el-col :span="12">
          <ko-form-item labelName="Subdomain" placeholder="e.g. web" itemType="input" v-model="form.subdomain" />
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span=12>
          <table style="width: 98%" class="tab-table">
            <tr>
              <th scope="col" width="93%" align="left">
                <label>nameservers</label>
              </th>
              <th align="left"></th>
            </tr>
            <tr v-for="row in form.dnsConfig.nameservers" v-bind:key="row.index">
              <td>
                <ko-form-item :withoutLabel="true" placeholder="e.g. 1.1.1.1" itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleNameserversDelete(row.index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleNameserversAdd">Add Nameserver</el-button>
              </td>
            </tr>
          </table>
        </el-col>
        <el-col :span=12>
          <table style="width: 98%" class="tab-table">
            <tr>
              <th scope="col" width="93%" align="left"><label>searches</label></th>
              <th align="left"></th>
            </tr>
            <tr v-for="row in form.dnsConfig.searches" v-bind:key="row.index">
              <td>
                <ko-form-item :withoutLabel="true" placeholder="e.g. mycompany.com" itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleSearchesDelete(row.index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleSearchesAdd">Add Search Domains</el-button>
              </td>
            </tr>
          </table>
        </el-col>
      </el-row>

      <div style="margin-top: 20px">
        <label>DNS Resolver Options
          <el-tooltip class="item" effect="dark" content="ProTip: Paste lines of key=value or key: value into any key field for easy bulk entry" placement="top-start">
            <i class="el-icon-question" />
          </el-tooltip>
        </label>
        <table style="width: 98%" class="tab-table">
          <tr>
            <th scope="col" width="48%" align="left"><label>key</label></th>
            <th scope="col" width="48%" align="left"><label>value</label></th>
            <th align="left"></th>
          </tr>
          <tr v-for="row in form.dnsConfig.options" v-bind:key="row.index">
            <td>
              <ko-form-item :withoutLabel="true" placeholder="e.g. foo" itemType="input" v-model="row.name" />
            </td>
            <td>
              <ko-form-item :withoutLabel="true" placeholder="e.g. bar" itemType="input" v-model="row.value" />
            </td>
            <td>
              <el-button type="text" style="font-size: 10px" @click="handleOptionsDelete(row.index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button @click="handleOptionsAdd">Add Search Domains</el-button>
            </td>
          </tr>
        </table>
      </div>

      <div style="margin-top: 20px">
        <label>Host Aliases
          <el-tooltip class="item" effect="dark" content="Additional /etc/hosts entries to be injected in the container." placement="top-start">
            <i class="el-icon-question" />
          </el-tooltip>
        </label>
        <table style="width: 98%" class="tab-table">
          <tr>
            <th scope="col" width="48%" align="left"><label>key</label></th>
            <th scope="col" width="48%" align="left"><label>value</label></th>
            <th align="left"></th>
          </tr>
          <tr v-for="(row, index) in form.hostAliases" v-bind:key="index">
            <td>
              <ko-form-item :withoutLabel="true" placeholder="e.g. 1.1.1.1" itemType="input" v-model="row.ip" />
            </td>
            <td>
              <ko-form-item :withoutLabel="true" placeholder="e.g. foo.com,bar.com" itemType="input" v-model="row.hostnames" />
            </td>
            <td>
              <el-button type="text" style="font-size: 10px" @click="handleAliasDelete(index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button @click="handleAliasAdd">Add Alians</el-button>
            </td>
          </tr>
        </table>
      </div>
    </ko-card>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoNetworking",
  components: { KoFormItem, KoCard },
  props: {
    networkingParentObj: Object,
  },
  data() {
    return {
      network_mode_list: [
        { label: "Normal", value: false },
        { label: "Host Network", value: true },
      ],
      dns_policy_list: [
        { label: "Default", value: "Default" },
        { label: "ClusterFirst", value: "ClusterFirst" },
        { label: "None", value: "None" },
      ],
      form: {
        subdomain: "",
        dnsPolicy: "",
        hostname: "",
        hostNetwork: "",
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
    transformation(parentFrom) {
      if (this.form.subdomain) {
        parentFrom.subdomain = this.form.subdomain
      }
      if (this.form.dnsPolicy) {
        parentFrom.dnsPolicy = this.form.dnsPolicy
      }
      if (this.form.hostname) {
        parentFrom.hostname = this.form.hostname
      }
      if (this.form.hostNetwork) {
        parentFrom.hostNetwork = true
      }
      if (this.form.hostAliases.length !== 0) {
        let aliases = []
        for (const item of this.form.hostAliases) {
          aliases.push({
            ip: item.ip,
            hostnames: item.hostnames.split(","),
          })
        }
        parentFrom.hostAliases = aliases
      }
      if (this.form.dnsConfig.nameservers.length !== 0 || this.form.dnsConfig.searches.length !== 0 || this.form.dnsConfig.options.length !== 0) {
        parentFrom.dnsConfig = {}
        if (this.form.dnsConfig.nameservers.length !== 0) {
          let nameservers = []
          for (const item of this.form.dnsConfig.nameservers) {
            nameservers.push(item.value)
          }
          parentFrom.dnsConfig.nameservers = nameservers
        }
        if (this.form.dnsConfig.searches.length !== 0) {
          let searches = []
          for (const item of this.form.dnsConfig.searches) {
            searches.push(item.value)
          }
          parentFrom.dnsConfig.searches = searches
        }
        if (this.form.dnsConfig.options.length !== 0) {
          parentFrom.dnsConfig.options = this.form.dnsConfig.options
        }
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
      if (this.networkingParentObj.hostNetwork) {
        this.hostNetwork = true
      }
      if (this.networkingParentObj.hostAliases) {
        for (const item of this.networkingParentObj.hostAliases) {
          this.form.hostAliases.push({
            ip: item.ip,
            hostnames: item.hostnames.join(","),
          })
        }
      }
      if (this.networkingParentObj.dnsConfig) {
        if (this.networkingParentObj.dnsConfig.nameservers) {
          for (const item of this.networkingParentObj.dnsConfig.nameservers) {
            this.form.dnsConfig.nameservers.push({ value: item })
          }
        }
        if (this.networkingParentObj.dnsConfig.searches) {
          for (const item of this.networkingParentObj.dnsConfig.searches) {
            this.form.dnsConfig.searches.push({ value: item })
          }
        }
        if (this.networkingParentObj.dnsConfig.options) {
          this.form.dnsConfig.options = this.networkingParentObj.dnsConfig.option
        }
      }
    }
  },
}
</script>