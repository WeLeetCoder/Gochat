<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
  :root {
      --font: Consolas,Monaco,Andale Mono,Ubuntu Mono,monospace, "Helvetica Neue",Helvetica,"PingFang SC","Hiragino Sans GB","Microsoft YaHei",Arial,sans-serif;
  }
  html {
    font-size: 12px;
    font-family: Consolas,Monaco,Andale Mono,Ubuntu Mono,monospace, "Helvetica Neue",Helvetica,"PingFang SC","Hiragino Sans GB","Microsoft YaHei",Arial,sans-serif;
  }
  a, button, div, input, select, span, strong, textarea {
    outline: 0;
  }
  * {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }
  pre, textarea{
    white-space: pre-wrap;
    word-break: break-word;
    font-family: var(--font);
  }
  ::-webkit-scrollbar {
    width: 8px;
  }
  ::-webkit-scrollbar-thumb {
    border-radius: 4px;
    background: rgba(179, 179, 179, 0.2);
  }
  ::-webkit-scrollbar-thumb:hover {
    background: rgba(179, 179, 179, 0.5);
  }

  ::-webkit-scrollbar-track {
    border-radius: 4px;
    background: none;
    padding: 2px;
  }
</style>
