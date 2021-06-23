<template>
    <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Groups' }">
        <el-row v-loading="loading">
            <el-col :span="4"><br/></el-col>
            <el-col :span="10">
                <div class="grid-content bg-purple-light">
                    <el-form ref="form" :model="form" label-width="150px" label-position="left">

                        <el-form-item :label="$t('commons.table.name')" prop="name" required>
                            <el-input v-model="form.name"></el-input>
                        </el-form-item>

                        <el-form-item :label="$t('business.user.role')" prop="roles">
                            <el-select v-model="form.roles" multiple placeholder="请选择">
                                <el-option
                                        v-for="item in roleOptions "
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                </el-option>
                            </el-select>
                        </el-form-item>

                        <el-form-item>
                            <div style="float: right">
                                <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                                <el-button type="primary" @click="onConfirm()">{{ $t("commons.button.confirm") }}
                                </el-button>
                            </div>
                        </el-form-item>
                    </el-form>
                </div>
            </el-col>
            <el-col :span="4"><br/></el-col>
        </el-row>
    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import {createGroup} from "@/api/groups"
    import {listRoles} from "@/api/roles"


    export default {
        name: "GroupCreate",
        components: {LayoutContent},
        data() {
            return {
                loading: false,
                roleOptions: [],
                form: {
                    name: "",
                    roles: [],
                },
            }
        },
        methods: {
            onConfirm() {
                const req = {
                    "apiVersion": "v1",
                    "kind": "User",
                    "name": this.form.name,
                    "roles": this.form.roles,
                }
                createGroup(req).then(() => {
                    this.$message({
                        type: "success",
                        message: this.$t("commons.msg.create_success")
                    })
                    this.$router.push({name: "Groups"})
                })
            },
            onCancel() {
                this.$router.push({name: "Groups"})
            },
        },
        created() {
            this.loading = true
            listRoles().then(d => {
                d.data.forEach(r => {
                    this.roleOptions.push({
                        label: r.name,
                        value: r.name,
                    })
                })
                this.loading = false
            })
        }
    }
</script>

<style scoped>
</style>
