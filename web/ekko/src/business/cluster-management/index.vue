<template xmlns:el-col="http://www.w3.org/1999/html">
    <layout-content>

        <el-row :gutter="10">
            <el-col :span="4" v-for="(item,index) in items" :key="index">
                <el-card class="box-card" shadow="hover">
                    <div v-if="!item.add">
                        <div slot="header" class="clearfix">
                            <el-row>
                                <el-col :span="18" style="font-size: 18px">
                                    {{item.name}}
                                </el-col>
                                <el-col :span="6">
                                    <el-button icon="el-icon-user" circle></el-button>
                                    <el-button icon="el-icon-delete" circle></el-button>
                                </el-col>
                            </el-row>
                        </div>
                        <div>
                            <el-row>
                                <el-col :span="4">
                                    <img width="38px" height="38px" src="~@/assets/kubernetes.png" alt="kubernetes.png">
                                </el-col>
                                <el-col :span="20">
                                    <table>
                                        <tr>
                                            <th style="width: 60%;text-align: left">版本</th>
                                            <td style="padding-left: 20px">{{item.status.version}}</td>
                                        </tr>
                                    </table>
                                </el-col>
                            </el-row>
                        </div>
                    </div>
                    <div v-if="item.add" @click="onCreate">
                        <i class="el-icon-plus"></i>
                    </div>
                </el-card>
            </el-col>
        </el-row>
    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import {listClusters} from "@/api/clusters"


    export default {
        name: "ClusterList",
        components: {LayoutContent},
        data() {
            return {
                items: [],
            }
        },
        methods: {
            onCreate() {
                this.$router.push({name: "ClusterCreate"})
            }
        },
        created() {
            listClusters().then(data => {
                this.items = data.data;
                this.items.push({
                    add: true,
                })
            })
        }
    }
</script>

<style scoped>

</style>