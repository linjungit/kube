<template>
    <a-layout>
        <!-- 固钉 不让header 进行移动 覆盖其他 -->
        <a-affix>
            <a-layout-header>
                <!-- 平台信息 -->
                <!-- float:left 居左的同时，会吧其他元素也放在同一行 -->
                <div style="float: left;">
                    <img style="height: 40px; margin-bottom: 10px;" :src="kubeLogo" />
                    <span style="font-size: 25px; padding: 0 50px 0 10px; font-weight: bold">kubeA</span>
                </div>
                <!-- 集群信息 -->
                <a-menu style="float:left;width: 250px;line-height:64px;" v-model:selectedKeys="selectedKeys1" theme="dark"
                    mode="horizontal">
                    <a-menu-item v-for="item in clusterList" :key="item">
                        {{ item }}
                    </a-menu-item>
                </a-menu>
                <!-- 用户信息 -->
                <div style="float:right">
                    <!-- 它将设置头像的高度为 40 像素，边框半径为 50%（使图像呈现圆形），并在图像右侧留出 10 像素的空白 -->
                    <img style="height: 40px; border-radius: 50%; margin-right: 10px;" :src="avator" />
                    <!-- 这里要注意要调整下拉框的样式 -->
                    <a-dropdown :overlayStyle="{ paddingTop: '22px' }">
                        <a>
                            Admin
                            <!-- 这个就是一个下拉图标 -->
                            <down-outlined />
                        </a>
                        <template #overlay>
                            <a-menu>
                                <a-menu-item>
                                    <a @click="Logout()">退出登录</a>
                                </a-menu-item>
                                <a-menu-item>
                                    <a>修改密码</a>
                                </a-menu-item>
                            </a-menu>
                        </template>
                    </a-dropdown>

                </div>
            </a-layout-header>
        </a-affix>
        <!-- 中部区域 高度永远都是窗口最大高度减去68px，因为68是 header 高度 -->
        <a-layout style="height: calc(100vh - 68px)">
            <!-- 侧边栏，核心 -->
            <a-layout-sider width="240" v-model:collapsed="collapsed" collapsible>
                <!-- selectedKeys表示点击选中的栏目 你点击会添加到这个数组里面,用于a-menu-item
                openKeys表示展开的栏目，用于a-sub-menu
                openChange事件监听SubMenu展开/关闭的回调  当我打开后会执行什么
                高度是100% 边框是0
            -->
            <a-menu
                :selectedKeys="selectedKeys2"     
                :openkeys="openKeys"
                @openChange="onOpenChange"
                mode="inline"
                :style="{height:'100%',boderRight:0}">
                <!-- 接下来拿到 routers 里面所有的路由信息生成对应的侧边栏 -->
                <template v-for="menu in routers" :key="menu"></template>
                    <!-- 处理无子路由的情况 -->
                    <a-menu-item
                    v-if="menu.children && menu.children.length ==1"
                    :index="menu.children[0].path"
                    :key="menu.children[0].path"
                    @click="routeChange()">
                    <template #icon>
                        <!-- 动态图标 -->
                        <component :is="menu.children[0].icon"></component>
                    </template>
                    <span>{{ menu.children[0].name }}</span>
                </a-menu-item>
                <!-- 处理子路由的情况，也就是父栏目 -->
                <a-sub-menu
                v-else-if="menu.children && menu.children.length > 1"
                :index="menu.path"
                :key="menu.path">
                <template #icon>
                    <component :is="menu.icon"/>
                </template>
                <template #title>
                    <span>
                        <span :class="[collapsed ? 'is-collapse' : '']">{{ menu.name }}</span>
                    </span>
                </template>
                </a-sub-menu>

            
            </a-menu>
            </a-layout-sider>
            <a-layout style="padding:0 24px">
                <!-- main 的部分 -->
                <!-- overflowY表示默认高度超出后，显示滚轮 -->
                <a-layout-content :style="{
                    background: 'rgb(31,30,30)',
                    padding: '24px',
                    margin: 0,
                    minHeight: '280px',
                    overflow: 'auto'
                }">
                    <router-view></router-view>
                </a-layout-content>
                <!-- footer 部分 -->
                <a-layout-footer style="text-align: center;">
                    2023 Created by addo Devops
                </a-layout-footer>
            </a-layout>
        </a-layout>
    </a-layout>
</template>

<script>
import { onMounted, ref } from 'vue';
// 这块是导入的图片
import kubeLogo from '@/assets/k8s-metrics.png';

import avator from '@/assets/avator.png';

export default ({
    setup() {
        const collapsed = ref(false)
        const selectedKeys1 = ref([])
        const clusterList = ref([
            'TST-1',
            'TST-2',
        ])
        // 拿到所有的路由信息
        const routers = ref([])
        const selectedKeys2 = ref([])
        const openkeys = ref([])
        // routers.value = router.options.routes 这里可以直接拿到router index里面的内容 export default router
        onMounted(() =>{
            routers.value = router.options.routes
        })
        // 通过useRouter方法获取路由配置以及当前页面的路由信息
        const  router = useRouter()

        function Logout() {
            //移出用户名 localStorage.removeItem 是 JavaScript 中用于从浏览器的本地存储
            localStorage.removeItem('username')
            //移出 token
            localStorage.removeItem('token')
            //跳转到 login 页面
            // router.push('login')
        }
        return {
            collapsed,
            kubeLogo,
            avator,
            selectedKeys1,
            clusterList,
            routers,
            selectedKeys2,
            openkeys,
            router,
            Logout,


        }
    }
})
</script>

<style>
.ant-layout-header {
    padding: 0 30px !important;
}

.ant-layout-footer {
    padding: 5px 50px !important;
    color: rgb(239, 239, 239);
}
</style>