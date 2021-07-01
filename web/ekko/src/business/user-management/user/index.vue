<template>
    <layout-content :header="$t('business.user.user_list')">
        <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                       :pagination-config="paginationConfig" @search="search">
            <template #header>
                <el-button-group>
                    <el-button v-has-permissions="[{resource:'users',verb:'create'}]" type="primary" size="small"
                               @click="onCreate">
                        {{ $t("commons.button.create") }}
                    </el-button>
                </el-button-group>
            </template>
            <el-table-column :label="$t('commons.table.name')" prop="name" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.name }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('business.user.nickname')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.nickName }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('business.user.email')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.email }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.creat_by')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.createdBy}}
                </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.createAt }}
                </template>
            </el-table-column>


            <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
        </complex-table>
    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import ComplexTable from "@/components/complex-table"
    import {searchUsers, deleteUser} from "@/api/users"
    import {checkPermissions} from "@/utils/permission";

    export default {
        name: "ClusterList",
        components: {ComplexTable, LayoutContent},
        data() {
            return {
                buttons: [
                    {
                        label: this.$t("commons.button.edit"),
                        icon: "el-icon-edit",
                        click: (row) => {
                            this.onEdit(row.name)
                        },
                        disabled: !checkPermissions({resource: "users", verb: "update"})
                    },
                    {
                        label: this.$t("commons.button.delete"),
                        icon: "el-icon-delete",
                        click: (row) => {
                            this.onDelete(row.name)
                        },
                        disabled: !checkPermissions({resource: "users", verb: "delete"})
                    },
                ],
                paginationConfig: {
                    currentPage: 1,
                    pageSize: 10,
                    total: 0,
                },
                searchConfig: {
                    quickPlaceholder: this.$t("commons.search.quickSearch"),
                    components: [
                        {field: "name", label: "姓名", component: "FuComplexInput"},
                        {
                            field: "status",
                            label: "状态",
                            component: "FuComplexSelect",
                            options: [{label: "启用", value: "Enable"}, {label: "禁用", value: "Disable"}],
                        },
                        // {
                        //     field: "name",
                        //     label: this.$t("commons.table.name"),
                        //     component: "FuComplexInput",
                        //     defaultOperator: "eq"
                        // },

                    ],
                },
                data: [],
                selects: []
            }
        },
        methods: {
            search(condition) {
                this.loading = true
                const {currentPage, pageSize} = this.paginationConfig
                searchUsers(currentPage, pageSize, condition).then(data => {
                    this.loading = false
                    this.data = data.data.items
                    this.paginationConfig.total = data.data.total
                })
            },
            onCreate() {
                this.$router.push({name: "UserCreate"})
            },
            onEdit(name) {
                this.$router.push({name: "UserEdit", params: {name: name}})
            },
            onDelete(name) {
                this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
                    confirmButtonText: this.$t("commons.button.confirm"),
                    cancelButtonText: this.$t("commons.button.cancel"),
                    type: 'warning'
                }).then(() => {
                    deleteUser(name).then(() => {
                        this.$message({
                            type: 'success',
                            message: this.$t("commons.msg.delete_success"),
                        });
                        this.search()
                    })
                });
            }

        },
        created() {
            this.search()
        }
    }
</script>

<style scoped>

</style>