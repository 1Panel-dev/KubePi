<template>
    <layout-content></layout-content>
</template>

<script>
    import LayoutContent from "@/components/layout/LayoutContent"
    import {listNamespace} from "@/api/namespace"

    export default {
        name: "NodeList",
        components: {LayoutContent},
        methods: {},
        mounted() {
            listNamespace().then(data => {
                console.log(data)
            })

        }
    }
</script>

<style scoped>

</style>