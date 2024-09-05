import {getConfigMap} from "@/api/configmaps"

export async function get_cluster_domain(clusterName){
    try{
       const cm = await getConfigMap(clusterName, "kube-system","coredns")

       const Corefile = cm.data["Corefile"];

       const Corefiles = Corefile.split(/[(\r\n)\r\n]+/); // 根据换行或者回车进行识别
    
       let clusterdomain=null;
       Corefiles.forEach((item, index) => {
        // 删除空项
        if (item && item != "") {
          if (item.trim().startsWith("kubernetes ")) {
            clusterdomain = item.trim().split(" ")[1];
          }
        }
       });
       return clusterdomain
      
    }catch(e){
        console.log(e)
        return null
    }
}
