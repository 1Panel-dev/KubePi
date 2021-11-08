<template>
      <div class="flag"
            :class="showFlag?'show':'hide'"
            @mouseenter="enter()"
            @mouseleave="leave()">{{showFlag?'Namespaced':'N'}}</div>

</template>

<script>

export default {
  name: "NamespacedTip",
  data() {
    return {
      showFlag: false,
    };
  },
  methods: {
    enter() {
      setTimeout(() => {
        this.showFlag = true;
      }, 100);
    },
    leave() {
      setTimeout(() => {
        this.showFlag = false;
      }, 200);
    },

  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    .flag {
        background: #e3e1e1;
        color: #666666;
        font-size: 12px;
        margin-left: 5px;
        display: inline-block;
        text-align: center;
        vertical-align: text-top;
        cursor: pointer;
        overflow: hidden;
        height: 20px;
        line-height: 20px;
        border-radius: 40%;
        -webkit-transition: all 0.3s ease;
        transition: all 0.3s ease;
        -webkit-transform: scale(0.8);
        -webkit-text-size-adjust: none;
    }
    .show {
        border-radius: 8px;
        width: 85px;
        /*margin-left: -3px;*/
    }
    .hide {
        border-radius: 40%;
        width: 20px;
    }
</style>
