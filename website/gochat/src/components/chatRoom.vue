<template>
    <div class="room-container">
        <section class="room-header">
            <Icon type="ios-chatbubbles" size="24" v-if="online"/>
            <Icon type="ios-chatbubbles-outline" size="24" v-else/>
            <span style="font-size: 1.4rem;">
                {{ courUser.Roomname }}
            </span> -
            <span style="font-size: 1.4rem;">
                {{ courUser.Username }}
            </span>
            - GoChat
        </section>
        <section class="room-body">
            <div class="room-chat-list">
                <ul>
                    <li class="chat-list-item" v-for="(chatList, index) in courUser.chatLists" :key="index" @click="selectSession(chatList)">
                        <div class="chat-list-icon-warp">
                            <div class="chat-picture"></div>
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
                        <div class="msg-item" v-for="(msgList, index) in msgTables"  :key="index" @contextmenu="usermenu" v-if="msgList.Content">
                            <div class="sender">
                                <div class="sender-picture"></div>
                            </div>
                            <div class="msg" >
                                <div class="msg-senderDetail">
                                    <span class="sender-name">
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
                        <div v-else class="msg-item">
                            <span class="sys-tip">{{msgList.user}} 上线啦！</span>
                        </div>
                    </div>
                    <div class="chat-msg-input">
                        <!-- 输入区域 -->
                        <!-- <pre class="msg-input" contenteditable="true"></pre> -->
                        <!--@keyup.enter="send"-->
                        <textarea  v-model="input" class="msg-input" placeholder="试试发送消息吧！" rows="5" ref="input" @keydown.enter="send"></textarea>
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
            online: false,
            input: '',
            courChat: {
                name: ' '
            },
            ctrl: false,
            msgLists: [],
            chatLists: [
              {
                name:"gochat",
              }
            ],
        }
    },
    methods: {
        send(e) {
            if (this.ctrl) {
                this.input += '\n';
                return
            }
            if (!this.input.trim()) {
              return ;
            }
            if (e !== undefined) {
                e.preventDefault();
            }
            let msg = {
                content: `<pre>${this.input}</pre>`.replace('\n', '<br>')
            };
            this.input = '';

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
        },
        usermenu(e) {
            e.preventDefault();
            console.log(e)
        }
    },
    computed: {
        msgTables() {
            if (this.$refs.msgDisplay) {
                this.$nextTick(() => {
                    let msgDisplay = this.$refs.msgDisplay;
                    let scrollHeight = msgDisplay.scrollHeight;
                    msgDisplay.scrollTop = scrollHeight;
                })
            }
            return this.msgLists
        },
    },
    mounted() {
        let info = getInfo();
        if (info) {
          this.courUser = info;
        }
        try {
          const wsUrl = 'ws://localhost:5000/chat/';
          let wsConn = new WebSocket(wsUrl + this.courUser.Token);

          wsConn.addEventListener('open', (e) => {
            this.wsConn = wsConn;
            // wsConn.send("login");
            this.online = true;
            console.log('ws connect opened');
          }, false);

          wsConn.addEventListener('message', (e) => {
            let response = JSON.parse(e.data);
            this.msgLists.push(response);
            console.log("receive msg:", response, );
          }, false);

          wsConn.addEventListener('error', (e) => {
            this.online = false;
            console.log('ws connect error');
          }, false);

          wsConn.addEventListener('close', (e) => {
            this.online = false;
            console.log('ws connect is closed');
          }, false);
        }catch (e) {
          this.online = false;
        }
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
                this.input += '\t';
            }
        });
    }
}
</script>
<style>
:root {
    --header-height: 48px;
    --icon-size: 40px;
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
    padding: 0 10px;
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
    font-size: 1.1rem;
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
    /*margin-bottom: 20px;*/
}
.msg-item {
    display: flex;
    padding: 20px 0;
    border-bottom: solid 1px rgba(255, 255, 255, .05);
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
    background: #666;
    -webkit-border-radius: 50%;
    -moz-border-radius: 50%;
    border-radius: 50%;
}
.msg-senderDetail {
    margin-bottom: 10px;
    font-size: 1.2rem;
}
.sender-time {
    font-size: .9rem;
    letter-spacing: 1px;
    color: hsla(0,0%,100%,.2);
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
    background: rgba(255, 255, 255, .2);
    -webkit-border-radius: 4px;
    -moz-border-radius: 4px;
    border-radius: 4px;
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
</style>
