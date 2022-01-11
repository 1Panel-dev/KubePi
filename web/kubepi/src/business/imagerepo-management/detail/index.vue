<template>
  <layout-content :header="$t('business.image_repos.images')" :back-to="{ name: 'ImageRepos' }">

    <el-alert
            :title="tip"
            type="info">
    </el-alert>
    <complex-table @search="search" :data="images" v-loading="loading">
      <el-table-column :label="$t('commons.table.name')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row }}
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {getRepo, listImagesByRepo} from "@/api/imagerepos"

export default {
  name: "ImageRepoDetail",
  components: { ComplexTable, LayoutContent },
  props: {
    repo: String,
  },
  data () {
    return {
      images: [],
      loading: false,
      repoObj: {},
      tip: "",
    }
  },
  methods: {
    search () {
      this.loading = true
      listImagesByRepo(this.repo).then(res => {
        this.images = res.data
      }).finally(() => {
        this.loading = false
      })
      getRepo(this.repo).then(res => {
        this.repoObj = res.data
        const repo = this.repoObj.repoName === ""? "":this.repoObj.repoName+"/"
        if (this.repoObj.type === 'Harbor') {
          this.tip =  this.$t('business.image_repos.push_image') + "docker push " + this.repoObj.downloadUrl + "/" + repo + "REPOSITORY[:TAG]"
        }else {
          this.tip =  this.$t('business.image_repos.push_image') + "docker push " + this.repoObj.downloadUrl  + "/REPOSITORY[:TAG]"
        }
      })
    }
  },
  created () {
    this.search()
  }
}
</script>

<style scoped>

</style>
