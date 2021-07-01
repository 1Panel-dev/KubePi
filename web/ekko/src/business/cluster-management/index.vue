<template xmlns:el-col="http://www.w3.org/1999/html">
    <layout-content :header="$t('business.cluster.cluster')">

        <div style="float: right">
            <el-switch
                    v-model="hiddenUnAccessCluster"
                    @change="onHiddenUnAccessClusterChange"
                    active-text="隐藏不可访问集群">
            </el-switch>
        </div>
        <br>
        <el-row>
            <el-col :span="4" v-for="(item,index) in items" :key="index">
                <el-card class="card_header box-card " style="margin-left:10px; margin-top:20px" shadow="hover">
                    <div>
                        <div slot="header" class="clearfix">
                            <b style="font-size: 20px">{{item.name}}</b>
                            <div style="float: right">
                                <el-dropdown trigger="click" @command="handleCommand">
                                    <el-button type="text" size="large" class="bottom-button">More</el-button>

                                    <el-dropdown-menu slot="dropdown">
                                        <el-dropdown-item v-has-permissions="{resource:'clusters',verb:'update'}"
                                                          :command="{name:item.name,action:'manage'}">管理集群
                                        </el-dropdown-item>
                                        <el-dropdown-item v-has-permissions="{resource:'clusters',verb:'delete'}"
                                                          :command="{name:item.name,action:'delete'}">删除
                                        </el-dropdown-item>
                                    </el-dropdown-menu>
                                </el-dropdown>

                            </div>
                            <hr style="border-color: gray"/>
                        </div>
                        <div>
                            <el-row>
                                <el-col :span="4" style="padding-top: 10px">
                                    <img height="48px" src="~@/assets/kubernetes.png" alt="kubernetes.png">
                                </el-col>
                                <el-col :span="20">
                                    <ul>
                                        <li style="list-style-type: none;margin: 5px">版本:{{item.status.version}}</li>
                                        <!--                                        <li style="list-style-type: none;margin: 5px">创建时间:{{item.createAt}}</li>-->
                                    </ul>
                                </el-col>
                            </el-row>
                        </div>
                        <div class="bottom clearfix">
                            <el-button type="text" size="large" class="bottom-button"
                                       v-has-permissions="{resource:'clusters',verb:'get'}"
                                       @click="onGotoDashboard(item.name)">Open
                            </el-button>
                        </div>

                    </div>

                </el-card>


            </el-col>
            <el-col :span="4">
                <el-card v-has-permissions="{resource:'clusters',verb:'update'}" class="card_header box-card "
                         style="margin-left:10px; margin-top:20px" shadow="hover">
                    <div @click="onCreate">
                        <i class="el-icon-plus"></i>
                    </div>
                </el-card>
            </el-col>
        </el-row>
    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import {listClusters, deleteCluster} from "@/api/clusters"


    export default {
        name: "ClusterList",
        components: {LayoutContent},
        data() {
            return {
                items: [],
                hiddenUnAccessCluster: false
            }
        },
        methods: {

            handleCommand(command) {
                const name = command.name
                switch (command.action) {
                    case "manage":
                        this.onDetail(name)
                        break
                    case "delete":
                        this.onDelete(name)
                        break
                }
            },

            onCreate() {
                this.$router.push({name: "ClusterCreate"})
            },
            onDetail(name) {
                this.$router.push({name: "ClusterMembers", params: {name: name}})
            },
            onDelete(name) {
                this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
                    confirmButtonText: this.$t("commons.button.confirm"),
                    cancelButtonText: this.$t("commons.button.cancel"),
                    type: 'warning'
                }).then(() => {
                    deleteCluster(name).then(() => {
                        this.$message({
                            type: 'success',
                            message: this.$t("commons.msg.delete_success"),
                        });
                        this.onVueCreated()
                    })
                });
            },
            onGotoDashboard(name) {
                window.open(`/dashboard?cluster=${name}`, "_blank")
            },
            onHiddenUnAccessClusterChange() {
                this.onVueCreated()
            },
            onVueCreated() {
                listClusters(!this.hiddenUnAccessCluster).then(data => {
                    this.items = data.data;
                })
            }
        },
        created() {
            this.onVueCreated()
        }
    }
</script>

<style scoped>
    .clearfix:before,
    .clearfix:after {
        display: table;
        content: "";
    }

    .clearfix:after {
        clear: both
    }

    .bottom {
        margin-top: 13px;
        line-height: 12px;
    }

    .bottom-button {
        padding: 0;
        float: right;
    }
</style>