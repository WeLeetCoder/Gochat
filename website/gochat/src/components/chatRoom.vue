<template>
    <div class="room-container">
        <section class="room-header">
            {{ courChat.name ? courChat.name : 'Username：' + courUser.Username}} - GoChat
        </section>
        <section class="room-body">
            <div class="room-chat-list">
                <ul>
                    <li class="chat-list-item" v-for="(chatList, index) in courUser.chatLists" :key="index" @click="selectSession(chatList)">
                        <div class="chat-list-icon-warp">
                            <div class="chat-picture">
                            </div>
                        </div>
                        <div class="chat-list-name">
                            {{ chatList.name }}
                        </div>
                    </li>
                </ul>
            </div>
            <div class="room-chat-msg" >
                <div class="room-chat-msg-wrap" v-if="courChat.name">
                    <div class="chat-msg-display" ref="msgDisplay">
                        <div v-for="(msgList, index) in msgTables" class="msg" :class="msgList.Sender == courUser.Id ? 'msg-me':'msg-other'" :key="index">

                            <div class="user-icon" v-if="msgList.Sender != courUser.Id">
                                <div class="chat-picture">
                                </div>
                            </div>
                            <div class="msg-content" v-else>
                                {{ msgList.Content}}
                            </div>

                            <div class="user-icon" v-if="msgList.Sender == courUser.Id">
                                <div class="chat-picture">
                                </div>
                            </div>
                            <div class="msg-content" v-else>
                                {{ msgList.Content}}
                            </div>

                        </div>
                        <!-- <div class="msg msg-me">
                            <div class="msg-content">
                                你是猴子请来的逗比吗？
                            </div>
                            <div class="user-icon">
                                <div class="chat-picture">
                                </div>
                            </div>
                        </div> -->
                    </div>
                    <div class="chat-msg-input">
                        <!-- 输入区域 -->
                        <!-- <pre class="msg-input" contenteditable="true"></pre> -->
                        <textarea  v-model="userEnterMsg" class="msg-input" placeholder="试试发送消息吧！" @keyup.enter="send"></textarea>
                        <div class="msg-sender">
                            <!-- 发送区域 -->
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
// import { request } from "../utils/ajax"
import { getInfo } from '../utils/store'

export default {
    data() {
        return {
            wsConn: null,
            courUser: {
              Username: "zhongqi",
              Token: "abcdefg",
            },
            userEnterMsg: '',
            courChat: {
                name: ' '
            },
            msgLists: [],
            chatLists: [
              {
                name:"gochat",
              }
            ],
        }
    },
    methods: {
        send() {
            if (!this.userEnterMsg) {
              return ;
            }
            let msg = {
                content: this.userEnterMsg
            };
            this.userEnterMsg = '';
            try {
              if (!this.wsConn) {
                console.log("ws connect failed, send failed");
                return
              }
              this.wsConn.send(JSON.stringify(msg));
            }catch (e) {
              console.log(e)
            }
        },
        selectSession(session) {
            this.courChat = session
        }
    },
    computed: {
        msgTables() {
            if (this.$refs.msgDisplay) {
                let msgDisplay = this.$refs.msgDisplay;
                let scrollHeight = msgDisplay.scrollHeight;
                msgDisplay.scrollTop = scrollHeight;
            }

            return this.msgLists
        }
    },
    mounted() {
        let info = getInfo();
        if (info) {
          this.courUser = info;
        }
        const wsUrl = 'ws://localhost:5000/chat/';
        let wsConn = new WebSocket(wsUrl + this.courUser.Token);

        wsConn.addEventListener('open', (e) => {
            this.wsConn = wsConn;
            wsConn.send("");
            console.log('ws connect opened');
        }, false);

        wsConn.addEventListener('message', (e) => {
             let response = JSON.parse(e.data);
             this.msgLists.push(response);
             console.log("receive msg:", response, );
        }, false);

        wsConn.addEventListener('error', (e) => {
            console.log('ws connect error');
        }, false);

        wsConn.addEventListener('close', (e) => {
            console.log('ws connect is closed');
        }, false);
    }
}
</script>
<style>
:root {
    --header-height: 48px;
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
    position: fixed;
    height: 100%;
    width: 100%;
    overflow: hidden;
    border-radius: 5px;
    box-shadow: 0 2px 10px 0 rgba(0,0,0,.2);
}
.room-header {
    height: var(--header-height);
    background: #36393f;
    color: white;
    text-align: center;
    font-size: 1.6rem;
    line-height: 48px;
}
.room-chat-list {
    width: 30%;
    min-width: 240px;
    max-width: 300px;
    height: 100%;
    background: #2f3136;
    border-right: solid 1px #2f3136;
    padding: 12px;
    overflow-y: auto;
}
.room-chat-msg-wrap {
    height: 100%;
}
.room-chat-msg {
    width: 80%;
    height: 100%;
    background: #16181d;
    color: white;
}
.room-body {
    display: flex;
    height: calc(100% - var(--header-height));
    position: relative;
}
.chat-list-item {
    height: 68px;
    border-bottom: solid 1px rgba(0, 0, 0, .2);
    display: flex;
}
.chat-list-icon-warp {
    height: 100%;
    width: 6rem;
    display: flex;
    align-items: center;
    justify-content: center;
}
.chat-picture {
    width: 48px;
    height: 48px;
    border-radius: 4px;
    background: url("../assets/icon.png") center;
    background-size: 100%;
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
    padding: 20px;
    padding-bottom: 0;
    overflow-y: auto;
    overflow-x: hidden;
}
.chat-msg-input {
    width: 100%;
    height: 30%;
    border-top: solid 1px #2f3136;
    display: flex;
    flex-direction: column;
    padding: 20px;
}
.msg-input {
    height: 70%;
    width: 100%;
    outline: none;
    border: 0;
    font-size: 1.3rem;
    overflow-y: auto;
    overflow-x: hidden;
    white-space: pre-wrap;
    word-break: normal;
    color: white;
    background: none;
    resize: none;
    margin-bottom: 20px;
}
.msg-sender {
    align-items: center;
    display: flex;
    justify-content: flex-end;
}
.msg {
    display: flex;
    margin-bottom: 20px;
    align-items: flex-start;
}
.msg-me {
    justify-content: flex-end;
}
.user-icon:nth-child(even) {
    margin-left: 10px;
}
.user-icon:nth-child(odd) {
    margin-right: 10px;
}
.msg-other {
    justify-content: flex-start;
}
.msg-content {
    min-height: 48px;
    background: #666;
    border-radius: 4px;
    padding: 10px;
    max-width: 70%;
    display: flex;
    align-items: center;
    font-size: 14px;
}
@media screen and (max-width: 480px) {
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
        flex-direction: row;
        height: 12%;
        position: fixed;
        bottom: 0;
    }
    .msg-input {
        height: 48px;
        width: 80%;
    }
}
</style>
