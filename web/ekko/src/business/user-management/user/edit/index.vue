<template>
    <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Users' }">
        <el-row v-loading="loading">
            <el-col :span="4"><br/></el-col>
            <el-col :span="10">
                <div class="grid-content bg-purple-light">
                    <el-form ref="form" :model="form" label-width="150px" label-position="left">

                        <el-form-item :label="$t('commons.table.name')" prop="name" required>
                            <el-input v-model="form.name"></el-input>
                        </el-form-item>


                        <el-form-item :label="$t('business.user.nickname')" prop="nickName" required>
                            <el-input v-model="form.nickName"></el-input>
                        </el-form-item>


                        <el-form-item :label="$t('business.user.email')" prop="email" required>
                            <el-input v-model="form.email"></el-input>
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
    import {updateUser, getUser} from "@/api/users"
    import {listRoles} from "@/api/roles"

    export default {
        name: "UserEdit",
        props: ["name"],
        components: {LayoutContent},
        data() {
            return {
                loading: false,
                user: {},
                roleOptions: [],
                form: {
                    name: "",
                    nickName: "",
                    email: "",
                    roles: []
                },
            }
        },
        methods: {
            onConfirm() {
                this.user.name = this.form.name
                this.user.spec.info.nickName = this.form.nickName
                this.user.spec.info.email = this.form.email
                this.user.roles = this.form.roles
                updateUser(this.name, this.user).then(() => {
                    this.$message({
                        type: "success",
                        message: this.$t("commons.msg.update_success")
                    })
                    this.$router.push({name: "Users"})
                })
            },
            onCancel() {
                this.$router.push({name: "Users"})
            },
            onCreated() {
                this.loading = true
                getUser(this.name).then(data => {
                    this.form.name = data.data.name;
                    this.form.nickName = data.data.spec.info.nickName;
                    this.form.email = data.data.spec.info.email;
                    this.form.roles = data.data.roles;
                    this.user = data.data;
                    listRoles().then(d => {
                        d.data.forEach(r => {
                            this.roleOptions.push({
                                label: r.name,
                                value: r.name,
                            })
                        })
                        this.loading = false
                    })
                })
            }
        },
        created() {
            this.onCreated()
        }
    }
</script>

<style scoped>
</style>
