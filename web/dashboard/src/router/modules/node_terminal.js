import terminals from "@/business/cluster/nodes/terminal"

const Terminal = {
    path: "/node_terminal?:cluster/:nodeName",
    hidden: true,
    component: terminals,
    name: "NodeTerminal",
    children: [
        {
            path: "/node_terminal",
            component: () => import("@/business/cluster/nodes/terminal"),
        },
    ]
}

export default Terminal
