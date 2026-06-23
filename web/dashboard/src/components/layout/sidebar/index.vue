<template>
    <div class="sidebar">
        <logo :collapse="isCollapse"/>
        <el-scrollbar wrap-class="scrollbar-wrapper">
            <el-menu
                    :default-active="activeMenu"
                    :collapse="isCollapse"
                    :collapse-transition="false"
                    :unique-opened="true"
                    mode="vertical"
                    active-text-color="#FFFFFF"
            >
                <sidebar-item v-for="route in permission_routes"
                              :key="route.path"
                              :item="route"
                              :base-path="route.path"/>
            </el-menu>
        </el-scrollbar>
        <div v-if="!isCollapse" class="sidebar-footer">
            <a class="sidebar-ad"
               href="https://1panel.cn"
               target="_blank"
               rel="noopener noreferrer">
                <span class="sidebar-ad-mark" aria-hidden="true">
                    <img :src="require('@/assets/1panel-menu-logo.svg')" class="sidebar-ad-logo" alt="">
                </span>
                <span class="sidebar-ad-content">
                    <span class="sidebar-ad-title">1Panel</span>
                    <span class="sidebar-ad-desc">{{ $t('commons.personal.one_panel_ad') }}</span>
                </span>
                <i class="el-icon-top-right sidebar-ad-icon"></i>
            </a>
            <button type="button" class="sidebar-copyright" @click="openApi">
                {{ $t('commons.personal.copy_right') }}
            </button>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex"
    import Logo from "@/components/layout/sidebar/Logo"
    import SidebarItem from "@/components/layout/sidebar/SidebarItem"

    export default {
        name: "Sidebar",
        components: {SidebarItem, Logo},
        data() {
            return {}
        },
        methods: {
          openApi(){
            window.open("https://www.fit2cloud.com/", "_blank");
          }
        },
        computed: {
            ...mapGetters([
                "permission_routes",
                "sidebar"
            ]),
            activeMenu() {
                const route = this.$route
                const {meta, path} = route
                if (meta.activeMenu) {
                    return meta.activeMenu
                }
                return path
            },
            isCollapse() {
                return !this.sidebar.opened
            },
        },
        created() {
        }
    }
</script>

<style lang="scss">
    @import "~@/styles/common/variables";

    @mixin sidebar-base-item {
        padding-left: 20px !important;
        border-radius: 2px;
        color: $menu-color;
    }

    @mixin menu-item {
        @include sidebar-base-item;
        line-height: $menu-height;
        height: $menu-height;
    }

    @mixin submenu-item {
        @include sidebar-base-item;
        line-height: $submenu-height;
        height: $submenu-height;
    }

    @mixin popper-submenu-item {
        @include sidebar-base-item;
        line-height: $submenu-height;
        height: $submenu-height;
    }

    @mixin menu-item-active {
        font-weight: 600;
        color: $menu-active-color;
        background-color: $menu-active-bg-color;
        &:hover {
            background-color: $menu-bg-color-hover;
        }
    }

    @mixin submenu-item-active {
        font-weight: 600;
        color: $submenu-active-color;
        background-color: $submenu-active-bg-color;
        &:hover {
            background-color: $submenu-bg-color-hover;
        }
    }

    @mixin menu-active-prefix {
        &:after {
            content: "";
            position: absolute;
            left: 0;
            top: 1px;
            bottom: 1px;
            width: $menu-active-prefix-width;
            background-color: $menu-active-prefix-color;
        }
    }

    .sidebar {
        height: calc(100% - 30px);
        display: flex;
        flex-direction: column;

        .el-scrollbar {
            box-sizing: border-box;
            padding: 10px 0;
            flex: 1;
            min-height: 0;
            height: auto;

            .el-scrollbar__bar {
                &.is-vertical {
                    right: 0;
                }

                &.is-horizontal {
                    display: none;
                }
            }

            .scrollbar-wrapper {
                height: 100%;
                overflow-x: hidden;
            }
        }

        .sidebar-footer {
            flex: 0 0 auto;
            padding: 8px 12px 14px;
            text-align: center;
        }

        .sidebar-ad {
            display: flex;
            align-items: center;
            gap: 10px;
            box-sizing: border-box;
            min-height: 58px;
            margin-bottom: 10px;
            padding: 10px;
            color: #e2e8f0;
            text-align: left;
            text-decoration: none;
            background: linear-gradient(135deg, rgba(37, 99, 235, 0.22), rgba(96, 165, 250, 0.08));
            border: 1px solid rgba(96, 165, 250, 0.24);
            border-radius: 8px;
            transition: border-color 0.2s ease, background 0.2s ease, transform 0.2s ease;

            &:hover,
            &:focus {
                color: #ffffff;
                background: linear-gradient(135deg, rgba(37, 99, 235, 0.34), rgba(96, 165, 250, 0.14));
                border-color: rgba(147, 197, 253, 0.46);
                transform: translateY(-1px);
                outline: none;
            }

            &:focus-visible {
                outline: 1px solid rgba(147, 197, 253, 0.62);
                outline-offset: 3px;
            }
        }

        .sidebar-ad-mark {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            flex: 0 0 28px;
            width: 28px;
            height: 28px;
            padding: 4px;
            box-sizing: border-box;
            background: $--color-primary;
            border: 1px solid rgba(255, 255, 255, 0.34);
            border-radius: 7px;
            box-shadow: 0 8px 18px rgba(37, 99, 235, 0.24);
        }

        .sidebar-ad-logo {
            display: block;
            width: 100%;
            height: 100%;
            object-fit: contain;
        }

        .sidebar-ad-content {
            display: flex;
            flex: 1;
            flex-direction: column;
            min-width: 0;
        }

        .sidebar-ad-title {
            color: #ffffff;
            font-size: 13px;
            font-weight: 600;
            line-height: 18px;
        }

        .sidebar-ad-desc {
            margin-top: 1px;
            color: rgba(226, 232, 240, 0.72);
            font-size: 12px;
            line-height: 18px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .sidebar-ad-icon {
            flex: 0 0 auto;
            color: rgba(226, 232, 240, 0.7);
            font-size: 14px;
        }

        .sidebar-copyright {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            max-width: 100%;
            min-height: 24px;
            padding: 0;
            border: 0;
            appearance: none;
            color: rgba(226, 232, 240, 0.62);
            background: transparent;
            font-family: inherit;
            font-size: 12px;
            line-height: 20px;
            white-space: nowrap;
            cursor: pointer;
            transition: color 0.2s ease;

            &:hover,
            &:focus {
                color: rgba(255, 255, 255, 0.88);
                background: transparent;
                outline: none;
            }

            &:focus-visible {
                outline: 1px solid rgba(148, 163, 184, 0.45);
                outline-offset: 3px;
            }
        }

        a {
            width: 100%;
            overflow: hidden;
        }

        .is-opened {
            .el-menu {
                background-color: $menu-open-bg-color;
            }
        }

        .el-menu {
            border: none;
            height: 100%;
            width: 100%;
            background-color: $menu-bg-color;

            .submenu-title-no-dropdown {
                @include menu-item;

                &:hover {
                    background-color: $menu-bg-color-hover;
                }

                &.is-active {
                    @include menu-item-active;
                    @include menu-active-prefix;
                }
            }

            .el-submenu {
                .el-submenu__title {
                    @include menu-item;

                    &:hover {
                        background-color: $menu-bg-color-hover;
                    }
                }

                &.is-active {
                    &:not(.is-opened) {
                        .el-submenu__title {
                            @include menu-active-prefix;
                        }
                    }

                    .el-submenu__title {
                        @include menu-item-active;

                        .sub-el-icon, span {
                            color: #FFF;
                        }
                    }
                }

                .el-menu-item {
                    @include submenu-item;

                    &:hover {
                        background-color: $submenu-bg-color-hover;
                    }

                    &.is-active {
                        @include submenu-item-active;
                        @include menu-active-prefix;
                    }
                }
            }

            .nest-menu, .el-submenu__title, .submenu-title-no-dropdown {
                span {
                    padding-left: 30px;
                }

                .sub-el-icon {
                    margin-right: 10px;

                    + span {
                        padding-left: 0;
                    }
                }
            }

            &.el-menu--collapse {
                .el-tooltip {
                    padding: 0 !important;
                    text-align: center;
                    line-height: $menu-height;
                }
                .submenu-title-no-dropdown, .el-submenu__title {
                    max-width: 60px;
                    text-align: center;

                    span {
                        display: none;
                    }

                    .sub-el-icon {
                        margin: 0;
                    }

                    .el-submenu__icon-arrow {
                        display: none;
                    }
                }
            }
        }
    }

    .sidebar-popper {
        .el-menu--popup {
            padding: 0;
        }

        & > .el-menu {
            display: block;
            background-color: $sidebar-bg-color;

            .sub-el-icon {
                margin-right: 12px;
                margin-left: -2px;
            }
        }

        .nest-menu .el-submenu > .el-submenu__title, .el-menu-item {
            &.is-active {
                color: $submenu-active-color;
                @include menu-active-prefix;
            }

            @include popper-submenu-item;

            span {
                padding-left: 30px;
            }

            .sub-el-icon {
                margin-right: 10px;

                + span {
                    padding-left: 0;
                }
            }

            &:hover {
                background-color: $submenu-bg-color-hover;
            }
        }

        .el-menu--popup {
            max-height: 100vh;
            overflow-y: auto;

            &::-webkit-scrollbar-track-piece {
                background: rgba(203, 213, 225, 0.18);
            }

            &::-webkit-scrollbar {
                width: 6px;
            }

            &::-webkit-scrollbar-thumb {
                background: rgba(148, 163, 184, 0.55);
                border-radius: 20px;
            }
        }
    }
</style>
