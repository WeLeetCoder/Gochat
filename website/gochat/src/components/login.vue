<template>
    <!-- 登录页面 -->
    <div class="login-continer">

        <section class="login-form">
            <div class="login-wrap">
                <div class="login-logo"></div>
                <div class="login-title">欢迎来到GoChat！</div>
                <div class="login-item">
                    <Input v-model="username" placeholder="输入您的昵称"  />
                </div>
                <div class="login-item">
                    <Input v-model="roomName" placeholder="输入您要进入的房间名称"  />
                </div>
                <div class="login-item" style="margin: 0; margin-top: 10px;">
                    <Button type="primary" long @click="getToken">进入房间</Button>
                </div>
            </div>
        </section>
        <div class="login-bg"></div>
    </div>
</template>
<script>
import {request} from '../utils/ajax'
import { setInfo } from '../utils/store'
export default {
    data() {
        return {
            username: '',
            roomName: '',
        }
    },
    methods: {
        getToken() {
            if (this.username && this.roomName) {
                request({
                    url: 'http://localhost:5000/join',
                    data: {
                        username: this.username,
                        roomname: this.roomName
                    },
                    method: "post"
                }).then(res => {
                    if (res.data.code > 1000) {
                        this.iv.Message.error(res.data.msg);
                        return
                    }
                    setInfo(res.data.data);
                    this.iv.Message.success('请稍后...');
                    setTimeout(() => {
                        this.$router.push('/chatroom')
                    }, 2000);
                    
                }).catch(error => {
                    this.iv.Message.error('发生错误意外错误');
                })
                return null;
            }
            this.iv.Message.error('用户名和房间名不能为空！');
        }
    }
}
</script>
<style>
html {
    width: 100%;
    height: 100%;
}
body {
    width: 100%;
    height: 100%;
}
.login-continer {
    height: 100%;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    position: fixed;
}
.login-logo {
    width: 80px;
    height: 80px;
    background: url("../assets/logo.png") center;
    background-size: 100%;
    margin-bottom: 10px;
}
.login-title {
    /* color: white; */
    font-size: 2rem;
    margin-bottom: 20px;
}
.login-form {
    width: 400px;
    /* height: 400px; */
    color: #72767d;
    background: white;
    box-shadow: 0 2px 10px 0 rgba(0,0,0,.2);
    padding: 40px;
    border-radius: 4px;
}
.login-wrap {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
}
.login-item {
    width: 100%;
    margin-bottom: 10px;
}
.login-bg {
    background: url("../assets/bg.png") #2f3136;
    top: 0;
    bottom: 0;
    right: 0;
    left: 0;
    z-index: -1;
    position: absolute;
}
</style>

