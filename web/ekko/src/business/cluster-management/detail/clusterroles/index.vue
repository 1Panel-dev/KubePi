<template>
    <layout-content>
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
                    {{ row.metadata.name }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.description')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{$t("business.cluster_role."+row.metadata.annotations["ekko-i18n"])}}
                </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.creat_by')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.metadata.annotations["created-by"] }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.metadata.annotations["created-at"] }}
                </template>
            </el-table-column>
            <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')" fix/>
        </complex-table>


        <el-dialog
                title="提示"
                :visible.sync="createDialogOpened"
                width="20%"
                center>
            <el-form :model="memberForm" label-position="left" label-width="120px">
                <el-form-item label="主体类型">
                    <el-select v-model="memberForm.subjectKind" @change="onSubjectKindChange">
                        <el-option value="User">
                            用户
                        </el-option>
                        <el-option value="Group">
                            用户组
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="主体名称">
                    <el-select v-model="memberForm.subjectName">
                        <el-option v-for="(item, index) in memberOptions" :key="index" :value="item.name">
                            {{item.name}}
                        </el-option>
                    </el-select>
                </el-form-item>

            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="createDialogOpened = false">取 消</el-button>
                <el-button type="primary" @click="onConfirm">确 定</el-button>
            </span>
        </el-dialog>
    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import ComplexTable from "@/components/complex-table"
    import {createClusterMember} from "@/api/clusters"
    import {listUsers} from "@/api/users"
    import {listGroups} from "@/api/groups"
    import {listClusterRoles} from "@/api/clusters";


    export default {
        name: "ClusterRoles",
        props: ["name"],
        components: {LayoutContent, ComplexTable},
        data() {
            return {
                createDialogOpened: false,
                memberOptions: [],
                memberForm: {
                    subjectKind: "",
                    subjectName: "",
                },
                buttons: [
                    {
                        label: this.$t("commons.button.edit"),
                        icon: "el-icon-edit",
                        click: (row) => {
                            this.onEdit(row.name)
                        }
                    },
                    {
                        label: this.$t("commons.button.delete"),
                        icon: "el-icon-delete",
                        click: (row) => {
                            this.onDelete(row.name)
                        }
                    },
                ],
                data: [],
            }
        },
        methods: {
            list() {
                this.loading = false
                listClusterRoles(this.name).then(data => {
                    this.loading = false
                    this.data = data.data
                })
            },
            onCreate() {
                this.memberForm.subjectKind = ""
                this.createDialogOpened = true
                this.onSubjectKindChange()
            },
            onDelete() {
                this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
                    confirmButtonText: this.$t("commons.button.confirm"),
                    cancelButtonText: this.$t("commons.button.cancel"),
                    type: 'warning'
                }).then(() => {
                });
            },
            onSubjectKindChange() {
                this.memberForm.subjectName = ""
                if (this.memberForm.subjectKind === 'Group') {
                    listGroups().then((data) => {
                        this.memberOptions = data.data;
                    })
                }
                if (this.memberForm.subjectKind === 'User') {
                    listUsers().then((data) => {
                        this.memberOptions = data.data;
                    })
                }
            },
            onConfirm() {
                createClusterMember(this.name, {
                    kind: this.memberForm.subjectKind,
                    name: this.memberForm.subjectName
                }).then(() => {
                    this.createDialogOpened = false
                    this.list()
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