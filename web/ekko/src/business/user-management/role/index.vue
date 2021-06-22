<template>
    <layout-content :header="$t('business.user.role_list')">
        <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                       :pagination-config="paginationConfig">
            <template #header>
                <el-button-group>
                    <el-button type="primary" size="small" @click="onCreate">
                        {{ $t("commons.button.create") }}
                    </el-button>
                </el-button-group>
            </template>
            <el-table-column type="selection" fix></el-table-column>
            <el-table-column :label="$t('commons.table.name')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.name }}
                </template>
            </el-table-column>

            <el-table-column :label="$t('commons.table.description')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ $t(row.description) }}
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
            <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')" fix/>
        </complex-table>
    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import ComplexTable from "@/components/complex-table"
    import {searchRoles, deleteRole} from "@/api/roles"


    export default {
        name: "Role",
        components: {LayoutContent, ComplexTable},
        data() {
            return {
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
                paginationConfig: {
                    currentPage: 1,
                    pageSize: 10,
                    total: 0,
                },
                searchConfig: {
                    quickPlaceholder: this.$t("commons.search.quickSearch"),
                    components: [
                        {
                            field: "name",
                            label: this.$t("commons.table.name"),
                            component: "FuComplexInput",
                            defaultOperator: "eq"
                        },
                        // { field: "created_at", label: this.$t("commons.table.create_time"), component: "FuComplexDateTime", valueFormat: "yyyy-MM-dd HH:mm:ss" },
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
                searchRoles(currentPage, pageSize, condition).then(data => {
                    this.loading = false
                    this.data = data.data.items
                    this.paginationConfig.total = data.data.total
                })
            },
            onCreate() {
                this.$router.push({name: "RoleCreate"})

            },
            onEdit(name) {
                this.$router.push({name: "RoleEdit", params: {name: name}})
            },
            onDelete(name) {
                this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
                    confirmButtonText: this.$t("commons.button.confirm"),
                    cancelButtonText: this.$t("commons.button.cancel"),
                    type: 'warning'
                }).then(() => {
                    deleteRole(name).then(() => {
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