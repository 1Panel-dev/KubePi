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
                    {{ row.name }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.kind')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.kind}}
                </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.createAt }}
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
                <el-form-item label="集群角色">
                    <el-select v-model="memberForm.clusterRoles" multiple placeholder="请选择">
                        <el-option
                                v-for="(item,index) in clusterRolesOptions"
                                :key="index"
                                :value="item.metadata.name">
                            {{item.metadata.name}}
                        </el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
    <el-button @click="createDialogOpened = false">取 消</el-button>
    <el-button type="primary" @click="onConfirm">确 定</el-button>
  </span>
        </el-dialog>


        <el-dialog
                title="提示"
                :visible.sync="editDialogOpened"
                width="20%"
                center>
            <el-form :model="editForm" label-position="left" label-width="120px">
                <el-form-item label="主体类型">
                    <el-select v-model="editForm.subjectKind" disabled @change="onSubjectKindChange">
                        <el-option value="User">
                            用户
                        </el-option>
                        <el-option value="Group">
                            用户组
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="主体名称">
                    <el-select v-model="editForm.subjectName" disabled>
                        <el-option v-for="(item, index) in memberOptions" :key="index" :value="item.name">
                            {{item.name}}
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="集群角色">
                    <el-select v-model="editForm.clusterRoles" multiple placeholder="请选择">
                        <el-option
                                v-for="(item,index) in clusterRolesOptions"
                                :key="index"
                                :value="item.metadata.name">
                            {{item.metadata.name}}
                        </el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
    <el-button @click="editDialogOpened = false">取 消</el-button>
    <el-button type="primary" @click="onEditConfirm">确 定</el-button>
  </span>
        </el-dialog>


    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import ComplexTable from "@/components/complex-table"
    import {createClusterMember, listClusterMembers} from "@/api/clusters"
    import {listUsers} from "@/api/users"
    import {listGroups} from "@/api/groups"
    import {listClusterRoles, deleteClusterMember, getClusterMember, updateClusterMember} from "@/api/clusters";


    export default {
        name: "ClusterMembers",
        props: ["name"],
        components: {LayoutContent, ComplexTable},
        data() {
            return {
                createDialogOpened: false,
                editDialogOpened: false,
                memberOptions: [],
                clusterRolesOptions: [],
                memberForm: {
                    subjectKind: "",
                    subjectName: "",
                    clusterRoles: []
                },
                editForm: {
                    subjectKind: "",
                    subjectName: "",
                    clusterRoles: []
                },
                buttons: [
                    {
                        label: this.$t("commons.button.edit"),
                        icon: "el-icon-edit",
                        click: (row) => {
                            this.onEdit(row)
                        }
                    },
                    {
                        label: this.$t("commons.button.delete"),
                        icon: "el-icon-delete",
                        click: (row) => {
                            this.onDelete(row)
                        }
                    },
                ],
                data: [],
            }
        },
        methods: {
            list() {
                this.loading = false
                listClusterMembers(this.name).then(data => {
                    this.loading = false
                    this.data = data.data
                })
            },
            onCreate() {
                this.memberForm.subjectKind = ""
                this.memberForm.clusterRoles = []
                this.createDialogOpened = true
                listClusterRoles(this.name).then(data => {
                    this.clusterRolesOptions = data.data
                })
                this.onSubjectKindChange()
            },

            onEdit(row) {
                listClusterRoles(this.name).then(data => {
                    this.clusterRolesOptions = data.data
                    getClusterMember(this.name, row.name, row.kind).then(data => {
                        this.editForm.subjectName = data.data.name
                        this.editForm.subjectKind = data.data.kind
                        this.editForm.clusterRoles = data.data.clusterRoles
                    })
                })
                this.editDialogOpened = true
            },
            onDelete(raw) {
                this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
                    confirmButtonText: this.$t("commons.button.confirm"),
                    cancelButtonText: this.$t("commons.button.cancel"),
                    type: 'warning'
                }).then(() => {
                    deleteClusterMember(this.name, raw.name, raw.kind).then(() => {
                        this.$message.success("删除成功")
                        this.list()
                    })
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
            onEditConfirm() {
                updateClusterMember(this.name, this.editForm.subjectName, {
                    name: this.editForm.subjectName,
                    kind: this.editForm.subjectKind,
                    clusterRoles: this.editForm.clusterRoles
                }).then(() => {
                    this.editDialogOpened = false
                    this.list()
                })
            },

            onConfirm() {
                createClusterMember(this.name, {
                    kind: this.memberForm.subjectKind,
                    name: this.memberForm.subjectName,
                    clusterRoles: this.memberForm.clusterRoles
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