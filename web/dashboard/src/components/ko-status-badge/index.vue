<template>
  <span v-if="hasValue" class="ko-status-badge" :class="statusClass">
    <span class="ko-status-badge__dot"></span>
    <span class="ko-status-badge__text">{{ value }}</span>
  </span>
</template>

<script>
const STATUS_TYPES = ["success", "warning", "danger", "info", "primary"]

export default {
  name: "KoStatusBadge",
  props: {
    value: {
      type: [String, Number],
      default: ""
    },
    type: {
      type: String,
      default: ""
    }
  },
  computed: {
    hasValue () {
      return this.value !== undefined && this.value !== null && `${this.value}` !== ""
    },
    normalizedType () {
      if (STATUS_TYPES.includes(this.type)) {
        return this.type
      }
      return this.getTypeByValue(this.value)
    },
    statusClass () {
      return `ko-status-badge--${this.normalizedType}`
    }
  },
  methods: {
    getTypeByValue (value) {
      const status = `${value}`.toLowerCase()
      if (["active", "running", "succeeded", "bound", "available", "ready"].includes(status)) {
        return "success"
      }
      if (["terminating", "pending", "waiting", "terminated", "notready", "schedulingdisabled"].includes(status)) {
        return "warning"
      }
      if (["failed", "error"].includes(status)) {
        return "danger"
      }
      return "info"
    }
  }
}
</script>

<style scoped>
.ko-status-badge {
  --ko-status-color: #475569;
  --ko-status-bg: #f8fafc;
  --ko-status-border: #cbd5e1;
  --ko-status-dot-shadow: rgba(100, 116, 139, 0.14);
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-width: 76px;
  height: 22px;
  padding: 0 10px;
  border: 1px solid var(--ko-status-border);
  border-radius: 999px;
  color: var(--ko-status-color);
  background: var(--ko-status-bg);
  font-size: 12px;
  font-weight: 600;
  line-height: 20px;
  box-sizing: border-box;
  vertical-align: middle;
}

.ko-status-badge__dot {
  width: 6px;
  height: 6px;
  flex: 0 0 6px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 0 3px var(--ko-status-dot-shadow);
}

.ko-status-badge__text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ko-status-badge--success {
  --ko-status-color: #15803d;
  --ko-status-bg: #f0fdf4;
  --ko-status-border: #bbf7d0;
  --ko-status-dot-shadow: rgba(34, 197, 94, 0.16);
}

.ko-status-badge--warning {
  --ko-status-color: #b45309;
  --ko-status-bg: #fffbeb;
  --ko-status-border: #fde68a;
  --ko-status-dot-shadow: rgba(245, 158, 11, 0.18);
}

.ko-status-badge--danger {
  --ko-status-color: #dc2626;
  --ko-status-bg: #fef2f2;
  --ko-status-border: #fecaca;
  --ko-status-dot-shadow: rgba(239, 68, 68, 0.16);
}

.ko-status-badge--info {
  --ko-status-color: #475569;
  --ko-status-bg: #f8fafc;
  --ko-status-border: #cbd5e1;
  --ko-status-dot-shadow: rgba(100, 116, 139, 0.14);
}

.ko-status-badge--primary {
  --ko-status-color: #2563eb;
  --ko-status-bg: #eff6ff;
  --ko-status-border: #bfdbfe;
  --ko-status-dot-shadow: rgba(37, 99, 235, 0.14);
}
</style>
