<template>
    <div class="errPage-container-1">
        <img src="@/assets/imgs/error_403.png" width="528" alt="403">
        <h2><span class="oops">Forbidden!</span>您暂无权限访问此资源!</h2>
        <p>{{authTime}}s 后自动返回首页</p>
        <el-button @click="$router.push('/')" round>返回首页</el-button>
    </div>
</template>


<script>
    export default {
        name: "Error403",
        data() {
            return {
                authTime: 0 /*倒计时 计数器*/,
            };
        },
        methods: {
            interval() {
                this.authTime = 5;
                let authTimetimer = setInterval(() => {
                    this.authTime--;
                    if (this.authTime <= 0) {
                        this.$router.push({path: "/"});
                        clearInterval(authTimetimer);
                    }
                }, 1000);
            },
        },
        created() {
            this.interval()
        }
    }
</script>

<style lang="scss" scoped>
    .errPage-container-1 {
        width: 80%;
        text-align: center;
        margin-top: 70px;
        color: #666666;

        .oops {
            font-size: 32px;
            font-weight: bold;
            color: #e27426;
            margin-right: 10px;
        }
    }
</style>