<template>
    <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Users' }">
        <el-row>
            <el-col :span="4"><br/></el-col>
            <el-col :span="10">
                <div class="grid-content bg-purple-light">
                    <el-form ref="form" :model="form" label-width="150px" label-position="left">

                        <el-form-item :label="$t('commons.table.name')" prop="name" required>
                            <el-input v-model="form.name"></el-input>
                        </el-form-item>


                        <el-form-item :label="$t('business.user.nickname')" prop="nickname" required>
                            <el-input v-model="form.nickname"></el-input>
                        </el-form-item>


                        <el-form-item :label="$t('business.user.email')" prop="email" required>
                            <el-input v-model="form.email"></el-input>
                        </el-form-item>


                        <el-form-item :label="$t('business.user.password')" prop="password" required>
                            <el-input type="password" v-model="form.password"></el-input>
                        </el-form-item>


                        <el-form-item :label="$t('business.user.confirm_password')" prop="confirmPassword" required>
                            <el-input type="password" v-model="form.confirmPassword"></el-input>
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
    import {createUser} from "@/api/users"

    export default {
        name: "UserCreate",
        components: {LayoutContent},
        data() {
            return {
                form: {
                    name: "",
                    nickname: "",
                    email: "",
                    password: "",
                    confirmPassword: "",
                },
            }
        },
        methods: {
            onConfirm() {
                const req = {
                    "apiVersion": "v1",
                    "kind": "User",
                    "name": this.form.name,
                    "spec": {
                        "info": {
                            "nickname": this.form.nickname,
                            "email": this.form.email,
                        },
                        "authenticate": {
                            "password": this.form.confirmPassword
                        }
                    }
                }
                createUser(req).then(() => {
                    this.$message({
                        type: "success",
                        message: this.$t("commons.msg.create_success")
                    })
                    this.$router.push({name: "Roles"})
                })
            },
            onCancel() {
                this.$router.push({name: "Users"})
            },
        },
        created() {
        }
    }
</script>

<style scoped>
</style>
