import {get, post, del, put} from "@/plugins/request";

const cornJobUrl = (cluster_name) => {
    return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/cronjobs`;
};
const cornJobWithNsUrl = (cluster_name, namespaces) => {
    return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/namespaces/${namespaces}/cronjobs`;
};

export function listCronJobs(cluster_name, search) {
    let url = cornJobUrl(cluster_name);
    if (search && search !== "") {
        url += "&fieldSelector=metadata.name=" + search;
    }
    return get(url);
}

export function getCronJobByName(cluster_name, namespace, cornjob) {
    return get(`${cornJobWithNsUrl(cluster_name, namespace)}/${cornjob}`);
}

export function deleteCronJob(cluster_name, cornjob) {
    return del(`${cornJobUrl(cluster_name)}/${cornjob}`);
}

export function createCronJob(cluster_name, cornjob) {
    return post(`${cornJobUrl(cluster_name)}/${cornjob}`);
}

export function updateCronJob(cluster_name, cornjob) {
    return put(`${cornJobWithNsUrl(cluster_name)}/${cornjob}`);
}
