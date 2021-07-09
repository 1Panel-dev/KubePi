<template>
    <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Roles' }">
        <el-row>
            <el-col :span="4"><br/></el-col>
            <el-col :span="10">
                <div class="grid-content bg-purple-light">
                    <el-form ref="form" :model="form" :rules="rules" label-width="150px" label-position="left">

                        <el-form-item :label="$t('commons.table.name')" prop="name" required>
                            <el-input v-model="form.name"></el-input>
                        </el-form-item>

                        <el-form-item :label="$t('commons.table.description')" prop="description">
                            <el-input v-model="form.description"></el-input>
                        </el-form-item>

                        <el-form-item>
                            <el-checkbox v-model="useTemplate">{{$t('business.user.base_on_exists_role')}}
                            </el-checkbox>
                        </el-form-item>
                        <el-form-item prop="template" v-if="useTemplate">
                            <el-select v-model="form.template" clearable>
                                <el-option
                                        v-for="item in roles"
                                        :key="item.name"
                                        :label="item.name"
                                        :value="item.name">
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
    import {listRoles, createRole} from "@/api/roles"
    import Rule from "@/utils/rules"

    export default {
        name: "RoleCreate",
        components: {LayoutContent},
        data() {
            return {
                loading: false,
                isSubmitGoing: false,
                roles: [],
                useTemplate: false,
                rules: {
                    name: [Rule.RequiredRule],
                },
                form: {
                    name: "",
                    description: "",
                    template: "",
                },
            }
        },
        methods: {
            onConfirm() {
                if (this.isSubmitGoing) {
                    return
                }
                this.isSubmitGoing = true
                const req = {
                    "apiVersion": "v1",
                    "kind": "Role",
                    "templateRef": this.form.template,
                    "name": this.form.name,
                    "description": this.form.description
                }
                createRole(req).then(() => {
                    this.$message({
                        type: "success",
                        message: this.$t("commons.msg.create_success")
                    })
                    this.$router.push({name: "Roles"})
                }).finally(() => {
                    this.isSubmitGoing = false
                })
            },
            onCancel() {
                this.$router.push({name: "Roles"})
            },
        },
        created() {
            listRoles().then(data => {
                this.roles = data.data
            })
        }
    }
</script>

<style scoped>
</style>
