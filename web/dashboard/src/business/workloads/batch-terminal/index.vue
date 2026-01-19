<template>
    <div style="background-color: #1f2224;height:100%;" >

        <el-row style="height:100%">
            <el-col :span="24" style="height:100%" >
                <el-tabs type="border-card" tabPosition="left" style="height:100%" v-model="activeName" >
                    <el-tab-pane :label="terminal.namespace + '.' + terminal.pod" :key="index" :name="terminal.namespace + '.' + terminal.pod"
                        v-for="(terminal, index) in terminals" :lazy="false">
                          <Terminal :params="terminal" :ref="setItemRef"/>
                    </el-tab-pane>
                </el-tabs>
            </el-col>
        </el-row>
    </div>
</template>
<script>
import  Terminal from '@/business/workloads/terminal'
export default {
    name: "BatchTerminal",
    components: {Terminal},
    data() {
        return {
            terminals: [],
            activeName: null,
            refList: [],
            refreshedMap:{}
        }
    },
    methods: {
        setItemRef(el) {
          this.refList.push(el)
        },
    },
    watch: {
        activeName(newVal) {
            this.$nextTick(() => {
                this.refList.forEach((item) => {
                    if(item.$parent.name === newVal) {
                      if(!this.refreshedMap[newVal]){  
                        item.getHeight()
                        item.reloadIframe()
                        this.refreshedMap[newVal]=true
                      }
                    }
                })
            })
        }
    },
    created() {
        this.terminals = JSON.parse(this.$route.query.terminals)
        if(this.terminals && this.terminals.length>0) {
            this.activeName = this.terminals[0].namespace + '.' + this.terminals[0].pod
        }
    },
}
</script>

<style lang="scss" scoped></style>

