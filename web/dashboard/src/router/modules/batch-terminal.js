import batchterminals from "@/business/workloads/batch-terminal"

const Terminal = {
    path: "/batch-terminal?:terminals",
    hidden: true,
    component: batchterminals,
    name: "BatchTerminal",
    children: [
        {
            path: "/batch-terminal",
            component: () => import("@/business/workloads/batch-terminal"),
        },
    ]
}

export default Terminal
