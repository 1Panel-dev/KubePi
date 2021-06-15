<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="Network Mode" clearable itemType="select" v-model="form.hostNetwork" :selections="network_mode_list" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="DNS Policy" clearable itemType="select" v-model="form.dnsPolicy" :selections="dns_policy_list" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Hostname" placeholder="e.g. web" clearable itemType="input" v-model="form.hostname" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Subdomain" placeholder="e.g. web" clearable itemType="input" v-model="form.subdomain" />
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span=12>
        <el-button @click="handleNameserversAdd">Add Nameserver</el-button>
        <el-table v-if="form.dnsConfig.nameservers.length !== 0" :data="form.dnsConfig.nameservers">
          <el-table-column min-width="80" label="Nameservers">
            <template v-slot:default="{row}">
              <ko-form-item :withoutLabel="true" placeholder="e.g. 1.1.1.1" clearable itemType="input" v-model="row.value" />
            </template>
          </el-table-column>
          <el-table-column width="120px">
            <template v-slot:default="{row}">
              <el-button type="text" style="font-size: 20px" @click="handleNameserversDelete(row)">REMOVE</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
      <el-col :span=12>
        <el-button @click="handleSearchesAdd">Add Search Domains</el-button>
        <el-table v-if="form.dnsConfig.searches.length !== 0" :data="form.dnsConfig.searches">
          <el-table-column min-width="80" label="Searches">
            <template v-slot:default="{row}">
              <ko-form-item :withoutLabel="true" placeholder="e.g. mycompany.com" clearable itemType="input" v-model="row.value" />
            </template>
          </el-table-column>
          <el-table-column width="120px">
            <template v-slot:default="{row}">
              <el-button type="text" style="font-size: 20px" @click="handleSearchesDelete(row)">REMOVE</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>

    <div style="margin-top: 20px">
      <label>DNS Resolver Options
        <el-tooltip class="item" effect="dark" content="ProTip: Paste lines of key=value or key: value into any key field for easy bulk entry" placement="top-start">
          <i class="el-icon-question" />
        </el-tooltip>
      </label>
      <div style="margin-top: 10px">
        <el-button @click="handleOptionsAdd">Add Search Domains</el-button>
      </div>
      <el-table v-if="form.dnsConfig.options.length !== 0" :data="form.dnsConfig.options">
        <el-table-column min-width="80" label="Key">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" placeholder="e.g. foo" clearable itemType="input" v-model="row.name" />
          </template>
        </el-table-column>
        <el-table-column min-width="80" label="Value">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" placeholder="e.g. bar" clearable itemType="input" v-model="row.value" />
          </template>
        </el-table-column>
        <el-table-column width="120px">
          <template v-slot:default="{row}">
            <el-button type="text" style="font-size: 20px" @click="handleOptionsDelete(row)">REMOVE</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div style="margin-top: 20px">
      <label>Host Aliases
        <el-tooltip class="item" effect="dark" content="Additional /etc/hosts entries to be injected in the container." placement="top-start">
          <i class="el-icon-question" />
        </el-tooltip>
      </label>
      <div style="margin-top: 10px">
        <el-button @click="handleAliasAdd">Add Alias</el-button>
      </div>
      <el-table v-if="form.hostAliases.length !== 0" :data="form.hostAliases">
        <el-table-column min-width="80" label="Key">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" placeholder="e.g. 1.1.1.1" clearable itemType="input" v-model="row.ip" />
          </template>
        </el-table-column>
        <el-table-column min-width="80" label="Value">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" placeholder="e.g. foo.com,bar.com" clearable itemType="input" v-model="row.hostname" />
          </template>
        </el-table-column>
        <el-table-column width="120px">
          <template v-slot:default="{row}">
            <el-button type="text" style="font-size: 20px" @click="handleAliasDelete(row)">REMOVE</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoNetworking",
  components: { KoFormItem },
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
        { label: "ClusterFirst", value: "Cluster First" },
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
    handleSearchesDelete(row) {
      for (let i = 0; i < this.form.dnsConfig.searches.length; i++) {
        if (this.form.dnsConfig.searches[i] === row) {
          this.form.dnsConfig.searches.splice(i, 1)
        }
      }
    },
    handleNameserversDelete(row) {
      for (let i = 0; i < this.form.dnsConfig.nameservers.length; i++) {
        if (this.form.dnsConfig.nameservers[i] === row) {
          this.form.dnsConfig.nameservers.splice(i, 1)
        }
      }
    },
    handleOptionsDelete(row) {
      for (let i = 0; i < this.form.dnsConfig.options.length; i++) {
        if (this.form.dnsConfig.options[i] === row) {
          this.form.dnsConfig.options.splice(i, 1)
        }
      }
    },
    handleAliasDelete(row) {
      for (let i = 0; i < this.form.hostAliases.length; i++) {
        if (this.form.hostAliases[i] === row) {
          this.form.hostAliases.splice(i, 1)
        }
      }
    },
    handleNameserversAdd() {
      var item = {
        value: "",
      }
      this.form.dnsConfig.nameservers.unshift(item)
    },
    handleSearchesAdd() {
      var item = {
        value: "",
      }
      this.form.dnsConfig.searches.unshift(item)
    },
    handleOptionsAdd() {
      var item = {
        name: "",
        value: "",
      }
      this.form.dnsConfig.options.unshift(item)
    },
    handleAliasAdd() {
      var item = {
        ip: "",
        hostname: "",
      }
      this.form.hostAliases.unshift(item)
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
  },
}
</script>