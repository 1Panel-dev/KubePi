<template>
    <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Users' }">
        <el-row v-loading="loading">
            <el-col :span="4"><br/></el-col>
            <el-col :span="10">
                <div class="grid-content bg-purple-light">
                    <el-form ref="form" :rules="rules" :model="form" label-width="150px" label-position="left">

                        <el-form-item :label="$t('business.user.username')" prop="name" required>
                            <el-input v-model="form.name" disabled></el-input>
                        </el-form-item>

                        <el-form-item :label="$t('business.user.nickname')" prop="nickname" required>
                            <el-input v-model="form.nickname"></el-input>
                        </el-form-item>


                        <el-form-item :label="$t('business.user.email')" prop="email" required>
                            <el-input v-model="form.email"></el-input>
                        </el-form-item>

                        <el-form-item :label="$t('business.user.role')" prop="roles">
                            <el-select v-model="form.roles"
                                       multiple
                                       style="width: 100%"
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
    import {updateUser, getUser} from "@/api/users"
    import {listRoles} from "@/api/roles"
    import Rules from "@/utils/rules"


    export default {
        name: "UserEdit",
        props: ["name"],
        components: {LayoutContent},
        data() {
            return {
                loading: false,
                isSubmitGoing: false,
                user: {},
                roleOptions: [],
                rules: {
                    name: [
                        Rules.RequiredRule
                    ],
                    email: [
                        Rules.RequiredRule,
                        Rules.EmailRule
                    ],
                    nickname: [
                        Rules.RequiredRule,
                    ],
                },
                form: {
                    name: "",
                    nickname: "",
                    email: "",
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
                this.loading = true

                this.user.name = this.form.name
                this.user.nickName = this.form.nickname
                this.user.email = this.form.email
                this.user.roles = this.form.roles
                updateUser(this.name, this.user).then(() => {
                    this.$message({
                        type: "success",
                        message: this.$t("commons.msg.update_success")
                    })
                    this.$router.push({name: "Users"})
                }).finally(() => {
                    this.isSubmitGoing = false
                    this.loading = false
                })
            },
            onCancel() {
                this.$router.push({name: "Users"})
            },
            onCreated() {
                this.loading = true
                getUser(this.name).then(data => {
                    this.form.name = data.data.name;
                    this.form.nickname = data.data.nickName;
                    this.form.email = data.data.email;
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
