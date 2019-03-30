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

  }
  html {
    font-size: 12px;
  }
  a, button, div, input, select, span, strong, textarea {
    outline: 0;
  }
  * {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }
  pre {
    white-space: pre-wrap;
    word-break: break-word;
    font-family: "Helvetica Neue",Helvetica,"PingFang SC","Hiragino Sans GB","Microsoft YaHei",Arial,sans-serif;
  }
  ::-webkit-scrollbar {
    width: 8px;
  }
  ::-webkit-scrollbar-thumb {
    border-radius: 4px;
    background: #2f3136;
  }
  ::-webkit-scrollbar-thumb:hover {
    background: #555;
  }

  ::-webkit-scrollbar-track {
    border-radius: 4px;
    background: none;
    padding: 2px;
  }
</style>
