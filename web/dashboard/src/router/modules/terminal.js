import terminals from "@/business/workloads/terminal"

const Terminal = {
    path: "/terminal?:cluster/:namespace/:pod/:container/:type",
    hidden: true,
    component: terminals,
    name: "Terminal",
    children: [
        {
            path: "/terminal",
            component: () => import("@/business/workloads/terminal"),
        },
    ]
}

export default Terminal