<template>
    <layout-content :header="$t('business.user.role_list')">
        <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                       :pagination-config="paginationConfig">
            <template #header>
                <el-button-group>
                    <el-button type="primary" size="small" @click="onImport">
                        {{ $t("commons.button.create") }}
                    </el-button>
                    <el-button type="primary" size="small" :disabled="selects.length===0">
                        {{ $t("commons.button.delete") }}
                    </el-button>
                </el-button-group>
            </template>
            <el-table-column type="selection" fix></el-table-column>
            <el-table-column :label="$t('commons.table.name')" min-width="100" fix>
                <template v-slot:default="{row}">
                    <el-link>{{ row.name }}</el-link>
                </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.create_time')" min-width="100" fix>
                <template v-slot:default="{row}">
                    {{ row.createAt }}
                </template>
            </el-table-column>
        </complex-table>
    </layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import ComplexTable from "@/components/complex-table"
    import {searchRoles} from "@/api/roles"


    export default {
        name: "Role",
        components: {LayoutContent, ComplexTable},
        data() {
            return {
                buttons: [
                    {
                        label: this.$t("commons.button.delete"),
                        icon: "el-icon-delete",
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
        },
        created() {
            this.search()
        }
    }
</script>

<style scoped>

</style>