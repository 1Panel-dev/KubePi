<template>
    <layout-content :header="$t('business.cluster.import')" :back-to="{ name: 'Clusters' }">
        <el-row>
            <el-col :span="4"><br/></el-col>
            <el-col :span="10">
                <div class="grid-content bg-purple-light">
                    <el-form :model="form" label-position="left" :rules="rules" label- width="60px">
                        <el-form-item :label="$t('commons.table.name')" prop="name">
                            <el-input v-model="form.name" clearable></el-input>
                        </el-form-item>
                        <el-divider content-position="center">连接设置</el-divider>
                        <el-form-item label="连接方式">
                            <el-radio v-model="form.direction" label="forward" @change="onDirectionChange">正向连接
                            </el-radio>
                            <el-radio v-model="form.direction" label="backward" @change="onDirectionChange">反向连接
                            </el-radio>
                        </el-form-item>
                        <div v-if="form.direction==='forward'">
                            <el-form-item label="API Server" prop="apiServer">
                                <el-input v-model="form.apiServer" clearable></el-input>
                            </el-form-item>

                            <el-form-item label="Insecure">
                                <el-switch v-model="form.apiServerInsecure"></el-switch>
                            </el-form-item>
                            <div v-if="!form.apiServerInsecure">
                                <el-form-item label="Ca Certificate" prop="caDataStr">
                                    <el-input type="textarea" v-model="form.caDataStr" clearable></el-input>
                                </el-form-item>
                            </div>


                            <el-form-item label="配置代理">
                                <el-switch v-model="form.proxyEnable"></el-switch>
                            </el-form-item>
                            <div v-if="form.proxyEnable">
                                <el-form-item label="Proxy Server" prop="proxyServer">
                                    <el-input v-model="form.proxyServer" clearable></el-input>
                                </el-form-item>
                                <el-form-item label="代理认证">
                                    <el-switch v-model="form.proxyAuthEnable"></el-switch>
                                </el-form-item>
                                <div v-if="form.proxyAuthEnable">
                                    <el-form-item label="代理用户名" prop="proxyUsername">
                                        <el-input v-model="form.proxyServer" clearable></el-input>
                                    </el-form-item>
                                    <el-form-item label="代理密码" prop="proxyPassword">
                                        <el-input v-model="form.proxyServer" clearable></el-input>
                                    </el-form-item>
                                </div>

                            </div>
                            <el-divider content-position="center">认证设置</el-divider>
                            <el-form-item label="认证方式">
                                <el-radio v-model="form.authenticationMode" label="bearer"
                                          @change="onAuthenticationModeChange">Bearer Token
                                </el-radio>
                                <el-radio v-model="form.authenticationMode" label="certificate"
                                          @change="onAuthenticationModeChange">证书
                                </el-radio>
                            </el-form-item>
                            <div v-if="form.authenticationMode==='bearer'">
                                <el-form-item label="Bearer Token" prop="token">
                                    <el-input type="textarea" v-model="form.token" clearable></el-input>
                                </el-form-item>
                            </div>
                            <div v-if="form.authenticationMode==='certificate'">
                                <el-form-item label="Certificate" prop="certDataStr">
                                    <el-input type="textarea" v-model="form.certDataStr" clearable></el-input>
                                </el-form-item>
                                <el-form-item label="Certificate Key" prop="keyDataStr">
                                    <el-input type="textarea" v-model="form.keyDataStr" clearable></el-input>
                                </el-form-item>
                            </div>
                        </div>
                        <el-form-item>
                            <div style="float: right">
                                <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                                <el-button type="primary" @click="onConfirm()" :disabled="isSubmitGoing">{{
                                    $t("commons.button.confirm") }}
                                </el-button>
                            </div>
                        </el-form-item>
                    </el-form>
                </div>
            </el-col>
        </el-row>
    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import {createCluster} from "@/api/clusters"

    export default {
        name: "ClusterCreate",
        components: {LayoutContent},
        data() {
            return {
                form: {
                    name: "",
                    direction: "forward",
                    apiServer: "",
                    apiServerInsecure: true,
                    authenticationMode: "bearer",
                    token: "",
                    proxyEnable: false,
                    proxyServer: "",
                    proxyAuthEnable: false,
                    proxyUsername: "",
                    proxyPassword: "",
                    certDataStr: "",
                    keyDataStr: "",
                    caDataStr: ""
                },
                rules: {},
                loading: false,
                isSubmitGoing: false,
            }
        },
        methods: {
            onDirectionChange() {
                if (this.form.direction === 'forward') {
                    this.form.apiServer = ""
                    this.form.authenticationMode = "bearer"
                    this.form.certificateData = ""
                    this.form.certificateKeyData = ""
                    this.form.token = ""
                }
            },
            onAuthenticationModeChange() {
                this.form.certificateData = ""
                this.form.certificateKeyData = ""
                this.form.token = ""
            },
            onCancel() {
                this.$router.push({name: "Clusters"})
            },
            onConfirm() {
                if (this.isSubmitGoing) {
                    return
                }
                this.loading = true
                const req = {
                    apiVersion: "v1",
                    kind: "Cluster",
                    name: this.form.name,
                    spec: {
                        connect: {
                            direction: this.form.direction
                        },
                        authentication: {
                            mode: this.form.authenticationMode
                        }

                    },
                    keyDataStr: "",
                    caDataStr: this.form.caDataStr,
                    certDataStr: ""
                }
                if (this.form.direction === 'forward') {
                    const forwardConfig = {}
                    forwardConfig.apiServer = this.form.apiServer
                    if (this.form.proxyEnable) {
                        forwardConfig.proxy.username = this.form.proxyUsername
                        forwardConfig.proxy.password = this.form.proxyPassword
                        forwardConfig.proxy.url = this.form.proxyServer
                    }
                    req.spec.connect.forward = forwardConfig
                }
                switch (this.form.authenticationMode) {
                    case "bearer":
                        req.spec.authentication.bearerToken = this.form.token
                        break
                    case "certificate":
                        req.certDataStr = this.form.certDataStr
                        req.keyDataStr = this.form.keyDataStr
                        break
                }
                console.log(req)
                createCluster(req).then(() => {
                    this.loading = false
                    this.isSubmitGoing = false
                    this.$message({
                        type: "success",
                        message: this.$t("commons.msg.create_success")
                    })
                    this.$router.push({"name": "Clusters"})
                })
            }
        },
    }
</script>

<style scoped>

</style>