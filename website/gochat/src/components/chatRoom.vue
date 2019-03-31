<template>
    <div class="room-container" ref="container">
        <section class="room-header" @mousedown="move" @mousemove="move" @mouseover="move">
            <ul class="ios-button">
                <li></li>
                <li></li>
                <li></li>
            </ul>

            <Icon type="ios-chatbubbles" size="24" v-show="online"/>
            <Icon type="ios-chatbubbles-outline" size="24" v-show="!online"/>
            <span style="font-size: 1.4rem;">
            {{ courUser.Roomname }}
            </span> -
            GoChat
        </section>
        <section class="room-body">
            <div class="room-chat-list">
                <ul>
                    <li class="chat-list-item">
                        <div class="chat-list-icon-warp">
                            <div class="room-icon"></div>
                        </div>
                        <div class="room-name">
                            <span style="font-size: 1.4rem; color: #fc4646; margin-right: 10px; max-width:35%;overflow: hidden;text-overflow:ellipsis; white-space: pre;">{{ courUser.Username }}</span>
                            <span style="font-size: 1.4rem;max-width:35%;overflow: hidden;text-overflow:ellipsis;white-space:pre;">@{{ roomName }}</span>
                        </div>
                    </li>
                    <li class="chat-list-item"
                        v-for="(chatList, index) in chatLists['user_list']"
                        :key="index"
                    >
                        <div class="chat-list-icon-warp">
                            <div class="chat-picture">
                                <span class="status" :class="online ? 'online' : 'offline'"
                                      v-if="courUser.Id === chatList.id"></span>
                            </div>
                        </div>
                        <div class="chat-list-name">
                            {{ chatList.name }}
                        </div>
                    </li>
                </ul>
            </div>
            <div class="room-chat-msg">
                <div class="room-chat-msg-wrap" v-if="courChat.name">
                    <div class="chat-msg-display" ref="msgDisplay">
                        <div class="msg-item" v-for="(msgList, index) in msgTables" :key="index" @contextmenu="usermenu"
                             v-if="msgList.Content">
                            <div class="sender">
                                <div class="sender-picture"></div>
                            </div>
                            <div class="msg">
                                <div class="msg-senderDetail">
                                    <span v-if="msgList.Sender === 'System'" style="color: #fc4646;">
                                            {{ msgList.SenderName }}
                                    </span>
                                    <span class="sender-name" v-else>
                                        {{ msgList.SenderName }}
                                    </span>

                                    <span class="sender-time">
                                        {{ msgList.Time }}
                                    </span>
                                </div>
                                <div class="msg-content" v-html="msgList.Content">
                                    <!--{{ msgList.Content }}-->
                                </div>
                            </div>
                        </div>
                        <div v-else class="msg-tip">
                            <span class="sys-tip">{{msgList.user}} {{ msgList.code === 999 ? "上线啦！": "下线了！" }}</span>
                        </div>
                    </div>
                    <div class="chat-msg-input">
                        <!-- 输入区域 -->
                        <!-- <pre class="msg-input" contenteditable="true"></pre> -->
                        <!--@keyup.enter="send"-->
                        <textarea v-model="input" class="msg-input" placeholder="试试发送消息吧！" rows="5" ref="input"
                                  @keydown.enter="send"></textarea>
                        <div class="msg-sender">
                            <!-- 发送区域 -->
                            <span class="tip">Ctrl+Enter 换行</span>
                            <Button type="info" @click="send()">发送</Button>
                        </div>
                    </div>
                </div>
                <div v-else style="font-size: 4rem; color: #666; text-align: center;margin-top: 200px;">
                    请选择一个会话吧！
                </div>
            </div>
        </section>
    </div>
</template>
<script>
    import {request} from "../utils/ajax";
    import {getInfo} from '../utils/store';
    import {SerPrefix} from "../config";

    export default {
        data() {
            return {
                wsConn: null,
                courUser: {
                    Username: '',
                    Token: '',
                    Roomname: '',
                    Id: '',
                },
                online: false,
                input: '',
                courChat: {
                    name: ' '
                },
                ctrl: false,
                msgLists: [],
                chatLists: {
                    "room_name": "GoChat",
                    "user_list": [
                        {
                            "id": "26523f1b2a8f7029e89f83c27082f11b",
                            "name": "dsfdsfdsfsdfsdfsdfsdfdsfsdfsdfdf"
                        },
                    ],
                },
            }
        },
        methods: {
            move(e) {
                this.$refs.container
            },
            send(e) {
                if (this.ctrl) {
                    this.input += '\n';
                    this.scrollBotton("input");
                    return
                }
                if (!this.input.trim()) {
                    return
                }
                if (e !== undefined) {
                    e.preventDefault();
                }
                let msg = {
                    content: this.input
                };
                this.input = '';

                try {
                    if (!this.wsConn) {
                        console.log("ws connect failed, send failed");

                        return
                    }
                    this.wsConn.send(JSON.stringify(msg));
                } catch (e) {
                    console.log(e)
                }
            },
            updateChat() {
                request({
                    url: `http://${SerPrefix}/users`,
                    data: {
                        "rname": this.courUser.Roomname,
                    }
                }).then((res) => {
                    if (!res.data.msg) {
                        this.chatLists = res.data.data
                    }

                })
            },
            usermenu(e) {
                e.preventDefault();
                console.log(e)
            },
            scrollBotton(tag) {
                this.$nextTick(() => {
                    let msgDisplay = this.$refs[tag];
                    let scrollHeight = msgDisplay.scrollHeight;
                    msgDisplay.scrollTop = scrollHeight;
                })

            }
        },
        computed: {
            msgTables() {
                if (this.$refs.msgDisplay) {
                    this.scrollBotton("msgDisplay")
                }
                return this.msgLists
            },
            roomName() {
                let roomname = this.courUser.Roomname;
                if (roomname) {
                    return roomname.toLocaleUpperCase();
                }
                return roomname
            }
        },
        mounted() {
            let info = getInfo();
            if (info) {
                this.courUser = info;
            }

            try {
                const wsUrl = `ws://${SerPrefix}/chat/`;
                let wsConn = new WebSocket(wsUrl + this.courUser.Token);

                wsConn.addEventListener('open', (e) => {
                    this.wsConn = wsConn;
                    this.online = true;
                    this.updateChat();
                }, false);

                wsConn.addEventListener('message', (e) => {
                    let response = JSON.parse(e.data);
                    if (response.user) {
                        this.updateChat();
                    }
                    this.msgLists.push(response);
                }, false);

                wsConn.addEventListener('error', (e) => {
                    this.online = false;
                    this.$router.push("/");
                    console.log('ws connect error');
                }, false);

                wsConn.addEventListener('close', (e) => {
                    this.online = false;
                    this.$router.push("/");
                    console.log('ws connect is closed');
                }, false);
            } catch (e) {
                this.online = false;
            }

            /*
                对 ctrl 和 tab 做处理
             */
            this.$refs.input.addEventListener('keyup', (e) => {
                if (e.key === 'Control') {
                    this.ctrl = false;
                }
                if (e.key === 'Tab') {
                    e.preventDefault();
                }
            });
            this.$refs.input.addEventListener('keydown', (e) => {
                if (e.key === 'Control') {
                    this.ctrl = true;
                }
                if (e.key === 'Tab') {
                    e.preventDefault();
                    this.input += '    ';
                }
            });
        }
    }
</script>
<style>
    :root {
        --header-height: 48px;
        --icon-size: 40px;
        --header-color: #212c44;
        font-size: 14px;
    }

    html {
        height: 100%;
    }

    body {
        height: 100%;
    }

    #app {
        background: url("../assets/bg.png") #d6d6d6;
        height: 100%;
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        overflow-y: auto;
    }

    .room-container {
        max-width: 1024px;
        max-height: 800px;
        margin: 0 auto;
        display: flex;
        flex-direction: column;
        position: fixed;
        height: 100%;
        width: 100%;
        overflow: hidden;
        border-radius: 5px;
        /*box-shadow: 0 2px 10px 0 rgba(0, 0, 0, .2);*/
        box-shadow: 2px 5px 15px 0px rgb(47, 49, 54);
    }

    .room-header {
        /*height: var(--header-height);*/
        background: #212c44;
        background: var(--header-color);
        color: white;
        text-align: center;
        font-size: 1.6rem;
        /*line-height: 48px;*/
        position: relative;
        user-select: none;
        cursor: default;
        padding: 5px;
    }

    .room-chat-list {
        width: 30%;
        min-width: 200px;
        max-width: 300px;
        height: 100%;
        background: #1c2333;
        padding: 12px;
        overflow-y: auto;
        overflow-x: hidden;
    }

    .room-chat-msg-wrap {
        height: 100%;
        padding: 5px;
    }

    .room-chat-msg {
        width: 80%;
        height: 100%;
        background: #262d3e;
        color: white;
    }

    .room-body {
        display: flex;
        height: 100%;
        position: relative;
    }

    .chat-list-item {
        height: 68px;
        border-bottom: solid 1px rgba(0, 0, 0, .2);
        display: flex;
        align-items: center;
    }

    .chat-list-icon-warp {
        height: 100%;
        min-width: 4rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .chat-picture {
        width: var(--icon-size);
        height: var(--icon-size);
        border-radius: 4px;
        background: url("../assets/icon.png") center no-repeat;
        background-size: cover;
        position: relative;
    }

    .chat-list-name {
        height: 100%;
        font-size: 1.5rem;
        display: flex;
        align-items: center;
        color: white;
    }

    .chat-msg-display {
        width: 100%;
        height: 70%;
        padding: 0 10px;
        overflow-y: auto;
        overflow-x: hidden;
    }

    .chat-msg-input {
        width: 100%;
        height: 30%;
        border-top: solid 1px #363a44;
        display: flex;
        flex-direction: column;
        -webkit-user-select: none;
        -moz-user-select: none;
        -ms-user-select: none;
        user-select: none;
        padding: 10px;
        padding-right: 0;
    }

    .msg-input {
        height: 70%;
        width: 100%;
        outline: none;
        border: 0;
        font-size: 1.1rem;
        overflow-y: auto;
        overflow-x: hidden;
        white-space: pre-wrap;
        word-break: normal;
        color: white;
        background: none;
        resize: none;
        padding-right: 5px;
        margin-bottom: 20px;
    }

    .msg-sender {
        align-items: center;
        display: flex;
        justify-content: flex-end;
        padding-right: 20px;
    }

    .msg {
        /*margin-bottom: 20px;*/
    }

    .msg-item {
        display: flex;
        padding: 10px;
        margin-bottom: 10px;
        /*border-bottom: solid 1px rgba(255, 255, 255, .05);*/
    }

    .msg-tip:first-child {
        padding-top: 10px;
    }

    .msg-tip {
        padding-bottom: 10px;
        display: flex;
        justify-content: center;
    }

    .msg-item:last-child {
        border-bottom: none;
    }

    .sender {
        padding-right: 10px;
    }

    .sender-picture {
        width: var(--icon-size);
        height: var(--icon-size);
        background: url("../assets/icon.jpg") #666 center;
        -webkit-background-size: 100%;
        background-size: 100%;
        -webkit-border-radius: 50%;
        -moz-border-radius: 50%;
        border-radius: 50%;
    }

    .msg-senderDetail {
        margin-bottom: 10px;
        font-size: 1.2rem;
        font-family: Consolas, Monaco, Andale Mono, Ubuntu Mono, monospace, "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", Arial, sans-serif;
    }

    .sender-time {
        font-size: .9rem;
        letter-spacing: 1px;
        color: hsla(0, 0%, 100%, .2);
    }

    .sender-name {
        font-size: 1.1rem;
        letter-spacing: 0;
        margin-right: 2px;

    }

    .msg-content {
        border-radius: 4px;
        display: flex;
        align-items: center;
        font-size: 1.1rem;
        line-height: 2rem;
    }

    .sys-tip {
        margin: 0 auto;
        padding: 5px 10px;
        background: rgba(173, 173, 173, .2);
        -webkit-border-radius: 4px;
        -moz-border-radius: 4px;
        border-radius: 4px;
        user-select: none;
    }

    .ios-button {
        list-style: none;
        position: absolute;
        left: 20px;
    }

    .ios-button > li {
        display: inline-block;
        border-radius: 50%;
        height: .9rem;
        width: .9rem;
        cursor: pointer;
    }

    .ios-button > li:nth-child(1) {
        background: #fc4646;
    }

    .ios-button > li:nth-child(2) {
        background: #fdb225;
    }

    .ios-button > li:nth-child(3) {
        background: #28c432;
    }

    .room-icon {
        min-width: 40px;
        max-width: 40px;
        max-height: 40px;
        min-height: 40px;
        border-radius: 4px;
        background: url("../assets/download.jpg") left no-repeat;
        background-size: cover;
        position: relative;
    }

    .room-name {
        min-width: 100%;
        display: flex;
        align-items: center;
        font-size: 2rem;
        color: white;
        font-family: var(--font);

    }

    .status {
        position: absolute;
        bottom: -4px;
        right: -2px;
        width: 1rem;
        height: 1rem;
        -webkit-border-radius: 50%;
        -moz-border-radius: 50%;
        border-radius: 50%;
        border: solid 2px #2f3136;
    }

    .online {
        background: #28c432;
    }

    .offline {
        background: #787878;
    }
    .tip {
        margin-right:10px;
        color: #6b7079d9;
    }
    @media screen and (max-width: 600px) {
        .tip {
            display: none;
        }
        .ios-button {
            display: none;
        }

        .room-container {
            height: 100%;
            max-width: none;
            max-height: none;
            border-radius: 0;
        }

        .room-chat-list {
            display: none;
            position: absolute;
        }

        .room-chat-msg {
            width: 100%;
        }

        .chat-msg-display {
            height: 88%;
        }

        .chat-msg-input {
            display: flex;
            align-items: center;
            flex-direction: row;
            height: 12%;
            position: fixed;
            bottom: 0;
            padding: 10px;
            justify-content: space-between;
        }

        .msg-input {
            /*height: 48px;*/
            margin-bottom: 0;
            margin-right: 10px;
            min-height: 24px;
        }
    }

    @media screen and (max-width: 768px) {
        .room-container {
            max-height: 100%;
        }
    }
</style>
