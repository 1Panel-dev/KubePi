<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.configuration.certificate')">
      <el-form label-position="top" ref="form" :model="form">
        <el-card v-for="(row,ind) in form.tls" v-bind:key="ind"
                 style="background-color: #292a2e;margin-top: 10px;position: relative">
          <div style="float: right;margin-top:-15px;position: relative;z-index: 1">
            <el-button type="text" v-if="form.tls.length > 1"
                       @click="removeTls(ind)">{{ $t("commons.button.delete") }}
            </el-button>
          </div>
          <div style="margin-top: 20px">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('business.configuration.certificate') + '- Secret Name'">
                  <el-select v-model="row.secretName" style="width: 100%">
                    <el-option label="Default Ingress Controller Certificate" value=""></el-option>
                    <el-option v-for="(secret,index) in secrets" :label="secret.metadata.name"
                               :value="secret.metadata.name" v-bind:key="index"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-row :gutter="20">
                  <div style="margin-top: 5px">
                    <el-col :span="24">
                      <span>Hosts</span>
                    </el-col>
                    <div v-for="(host,index) in row.hosts" v-bind:key="index">
                      <el-col :span="20">
                        <span v-if="false">{{ host }}</span>
                        <el-input v-model="row.hosts[index]" style="margin-top: 5px"></el-input>
                      </el-col>
                      <el-col :span="2">
                        <el-button type="text" @click="removeHost(row,index)">{{
                            $t("commons.button.delete")
                          }}
                        </el-button>
                      </el-col>
                    </div>
                    <div>
                      <el-col :span="24">
                        <el-button style="margin-top: 10px" @click="addHost(row)">{{ $t("commons.button.add") }} Host
                        </el-button>
                      </el-col>
                    </div>
                  </div>
                </el-row>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-form>
      <div style="margin-top: 10px">
        <el-button @click="addTls">{{ $t("commons.button.add") }}{{$t('business.configuration.certificate')}}</el-button>
      </div>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
import {listSecretsWithNs} from "@/api/secrets"

export default {
  name: "KoIngressTls",
  components: { KoCard },
  props: {
    cluster: String,
    namespace: String,
    tlsArray: Array
  },
  data () {
    return {
      form: {
        tls: []
      },
      secrets: []
    }
  },
  methods: {
    removeTls (index) {
      this.form.tls.splice(index, 1)
    },
    addTls () {
      this.form.tls.push({
        secretName: "",
        hosts: []
      })
    },
    addHost (row) {
      row.hosts.push("")
    },
    removeHost (row, index) {
      row.hosts.splice(index, 1)
    },
    getSecrets () {
      listSecretsWithNs(this.cluster, this.namespace).then(res => {
        this.secrets = res.items
      })
    }
  },
  watch: {
    namespace: function (oldValue, newValue) {
      if (oldValue !== newValue) {
        for (let i = 0; i < this.form.tls.length; i++) {
          this.form.tls[i].secretName = ""
        }
        this.getSecrets()
      }
    }
  },
  created () {
    this.form.tls = this.tlsArray ? this.tlsArray : [{
      secretName: "",
      hosts: []
    }]
    this.$emit("update:tlsArray", this.form.tls)
    this.getSecrets()
  }
}
</script>

<style scoped>

</style>
