<template>
    <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Groups' }">
        <el-row v-loading="loading">
            <el-col :span="4"><br/></el-col>
            <el-col :span="10">
                <div class="grid-content bg-purple-light">
                    <el-form ref="form" :model="form" label-width="150px" label-position="left">

                        <el-form-item :label="$t('commons.table.name')" prop="name" required>
                            <el-input v-model="form.name" disabled></el-input>
                        </el-form-item>

                        <el-form-item :label="$t('business.user.role')" prop="roles">
                            <el-select v-model="form.roles" multiple
                                       :placeholder="$t('commons.form.select_placeholder')">
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
    import {getGroup, updateGroup} from "@/api/groups"
    import {listRoles} from "@/api/roles"

    export default {
        name: "GroupEdit",
        props: ["name"],
        components: {LayoutContent},
        data() {
            return {
                loading: false,
                isSubmitGoing: false,
                user: {},
                roleOptions: [],
                group: {},
                form: {
                    name: "",
                    roles: []
                },
            }
        },
        methods: {
            onConfirm() {
                if (this.isSubmitGoing) {
                    return
                }
                this.isSubmitGoing = true

                this.group.name = this.form.name
                this.group.roles = this.form.roles
                updateGroup(this.name, this.group).then(() => {
                    this.$message({
                        type: "success",
                        message: this.$t("commons.msg.update_success")
                    })
                    this.$router.push({name: "Groups"})
                }).finally(() => {
                    this.isSubmitGoing = false
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
                getGroup(this.name).then(data => {
                    this.form.name = data.data.name
                    this.form.roles = data.data.roles
                    this.group = data.data;
                    this.loading = false
                })
            })
        }
    }
</script>

<style scoped>
</style>
