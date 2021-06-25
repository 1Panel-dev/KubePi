<template>
    <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Groups' }">
        <complex-table :data="data">
            <template #header>
                <el-button-group>
                    <el-button type="primary" size="small" @click="onCreate">
                        {{ $t("commons.button.create") }}
                    </el-button>
                </el-button-group>
                <br/>
            </template>
            <el-table-column :label="$t('commons.table.name')" min-width="100" fix>
                <template v-slot:default="{row}">
                    <el-link>{{ row.name }}</el-link>
                </template>
            </el-table-column>

            <el-table-column :label="$t('commons.table.creat_by')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.createBy}}
                </template>
            </el-table-column>


            <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.createAt }}
                </template>
            </el-table-column>
            <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
        </complex-table>

        <el-dialog
                title="提示"
                :visible.sync="createDialogOpened"
                width="20%"
                center z-index="20">
            <el-form :model="form" label-position="left" label- width="60px">
                <el-form-item label="Username">
                    <el-select v-model="form.username">
                        <el-option v-for="(item,index) in userOptions" :key="index"
                                   :value="item.name">
                            {{item.name}}
                        </el-option>
                    </el-select>
                </el-form-item>
            </el-form>


            <span slot="footer" class="dialog-footer">
                <el-button @click="createDialogOpened=false">取 消</el-button>
                <el-button type="primary" @click="onConfirm">确 定</el-button>
            </span>
        </el-dialog>

    </layout-content>
</template>

<script>
    import ComplexTable from "@/components/complex-table"
    import LayoutContent from "@/components/layout/LayoutContent"
    import {listGroupMembers, createGroupMember, deleteGroupMember} from "@/api/groups"
    import {listUsers} from "@/api/users"

    export default {
        name: "GroupBinding",
        props: ["name"],
        components: {LayoutContent, ComplexTable},
        data() {
            return {
                loading: false,
                createDialogOpened: false,
                userOptions: [],
                data: [],
                form: {
                    username: "",
                },
                buttons: [
                    {
                        label: this.$t("commons.button.delete"),
                        icon: "el-icon-delete",
                        click: (row) => {
                            this.onDelete(row.name)
                        }
                    },
                ],
            }
        },
        methods: {
            onCancel() {
                this.$router.push({name: "Groups"})
            },
            onConfirm() {
                createGroupMember(this.name, {
                    name: this.form.username,
                    groupName: this.name
                }).then(() => {
                    this.createDialogOpened = false
                    this.list()
                })
            },
            onDelete(name) {
                this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
                    confirmButtonText: this.$t("commons.button.confirm"),
                    cancelButtonText: this.$t("commons.button.cancel"),
                    type: 'warning'
                }).then(() => {
                    deleteGroupMember(this.name, name).then(() => {
                        this.$message({
                            type: 'success',
                            message: this.$t("commons.msg.delete_success"),
                        });
                        this.list()
                    })
                });
            },
            list() {
                listGroupMembers(this.name).then(data => {
                    this.data = data.data
                    console.log(this.data)
                })
            },
            onCreate() {
                this.createDialogOpened = true
                listUsers().then(data => {
                    this.userOptions = data.data;
                })
            }
        },
        created() {
            this.list()
        }
    }
</script>

<style scoped>
</style>
