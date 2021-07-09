<template>
    <layout-content :header="$t('business.user.user_group_list')">
        <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                       :pagination-config="paginationConfig">
            <template #header>
                <el-button-group>
                    <el-button v-has-permissions="[{resource:'groups',verb:'create'}]" type="primary" size="small"
                               @click="onCreate">
                        {{ $t("commons.button.create") }}
                    </el-button>
                </el-button-group>
            </template>
            <el-table-column :label="$t('commons.table.name')" min-width="100" fix>
                <template v-slot:default="{row}">

                    <el-link @click="onDetail(row.name)">{{ row.name }}</el-link>
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
    import {searchGroups, deleteGroup} from "@/api/groups"
    import {checkPermissions} from "@/utils/permission";


    export default {
        name: "Group",
        components: {LayoutContent, ComplexTable},
        data() {
            return {
                buttons: [
                    {
                        label: this.$t("commons.button.edit"),
                        icon: "el-icon-edit",
                        click: (row) => {
                            this.onEdit(row.name)
                        },
                        disabled: !checkPermissions({resource: "groups", verb: "update"})
                    },
                    {
                        label: this.$t("commons.button.delete"),
                        icon: "el-icon-delete",
                        click: (row) => {
                            this.onDelete(row.name)
                        },
                        disabled: !checkPermissions({resource: "groups", verb: "delete"})
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
            onCreate() {
                this.$router.push({name: "GroupCreate"})
            },
            onEdit(name) {
                this.$router.push({name: "GroupEdit", params: {name: name}})
            },
            onDetail(name) {
                this.$router.push({name: "GroupBinding", params: {name: name}})
            },
            onDelete(name) {
                this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
                    confirmButtonText: this.$t("commons.button.confirm"),
                    cancelButtonText: this.$t("commons.button.cancel"),
                    type: 'warning'
                }).then(() => {
                    deleteGroup(name).then(() => {
                        this.$message({
                            type: 'success',
                            message: this.$t("commons.msg.delete_success"),
                        });
                        this.search()
                    })
                });
            },
            search(condition) {
                this.loading = true
                const {currentPage, pageSize} = this.paginationConfig
                searchGroups(currentPage, pageSize, condition).then(data => {
                    this.loading = false
                    this.data = data.data.items
                    this.paginationConfig.total = data.data.total
                })
            },
        },
        created() {
            this.search()
        }
    }
</script>

<style scoped>

</style>