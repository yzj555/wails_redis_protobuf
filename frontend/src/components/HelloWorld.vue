<template>
  <div id="outerLayer">
    <v-app class="app">
      <div class="centerArea">
        <div class="queryArea">
          <v-row>
            <v-col>
              <v-text-field label="Redis-Key" v-model="keyWords" density="compact" @click:append="searchKey"
                append-icon="mdi-magnify"></v-text-field>
              <!-- <v-autocomplete label="Redis-Key" v-model="queryKey" :items="keysResult" density="compact"
                @click:append="searchKey" append-icon="mdi-magnify"></v-autocomplete> -->
            </v-col>
            <v-col>
              <v-text-field label="hash-key" v-model="queryField" density="compact"></v-text-field>
            </v-col>
            <v-col>
              <v-autocomplete v-model="parseMode" :items="parseModeList" item-title="text" item-value="value"
                density="compact"></v-autocomplete>
            </v-col>
            <v-col>
              <v-autocomplete v-model="queryMsg" :items="msgNameList" density="compact"
                label="消息名称(为空则返回源数据)"></v-autocomplete>
            </v-col>
            <v-col>
              <v-btn depressed small color="primary" @click="getDataByKeyAndMsg">
                解析
              </v-btn>
            </v-col>
          </v-row>
        </div>
        <div class="showArea">
          <!-- <pre v-html="highlightJSON(queryResult)" class="json-highlight">
                    </pre> -->
          <!-- {{ queryResult }} -->
          <div class="showKeys">
            <v-list density="compact">
              <v-list-item v-for="(item, i) in keysResult" :key="i" :value="item" active-color="green"
                @click="selectKey(item)">
                <v-list-item-title v-text="item"></v-list-item-title>
              </v-list-item>
            </v-list>
          </div>
          <!-- <pre id="json-renderer" class="showRes json-highlight" v-show="Object.keys(queryResult).length != 0"></pre> -->
          <pre id="json-renderer" class="showRes json-highlight"></pre>
        </div>
      </div>
    </v-app>
  </div>
</template>

<script>
import { Greet } from '../../wailsjs/go/main/App'
export default {
  name: "HelloWorld",
  data() {
    return {
      parseModeList: [
        { text: 'ProtoBuf', value: 1 },
        { text: '源数据', value: 2 },
        { text: 'Msgpack', value: 3 },
      ],
      beforeUnloadTime: 0,
      msgNameList: [],
      keyWords: "*",
      keysResult: [],
      queryKey: "",
      queryField: "",
      parseMode: 1, //默认pb
      queryMsg: "",
      queryResult: {},
    }
  },
  mounted() {
    this.getAllMsg()
    window.addEventListener('beforeunload', () => { this.beforeUnloadTime = Date.now() });
    window.addEventListener('unload', this.handleBeforeUnload);
  },
  methods: {
    handleBeforeUnload() {
      const gap = Date.now() - this.beforeUnloadTime;
      if (gap < 2) { // 阈值需实测调整
        console.log("页面关闭");
        this.requestSingle('/closeServer', {})
      } else {
        console.log("页面刷新");
      }
    },
    async getAllMsg() {
      Greet('HandleAllMsgName', JSON.stringify({})).then(result => {
        console.log(result)
        this.msgNameList = result.data
      })
    },

    async searchKey() {
      if (this.keyWords == "") {
        this.keyWords = "*"
      }
      Greet('HandleAllKeys', JSON.stringify({ likeKey: this.keyWords })).then(result => {
        console.log(result)
        this.keysResult = result.data
      })
    },
    selectKey(key) {
      this.queryKey = key
      this.getDataByKeyAndMsg()
    },
    async getDataByKeyAndMsg() {
      var params = {
        key: this.queryKey,
        field: this.queryField,
        parseMode: this.parseMode + "",
        msg: this.queryMsg,
      }
      Greet('HandleDataByKey', JSON.stringify(params)).then(result => {
        this.queryResult = result.data
        console.log(result)
        try {
          const viewer = new JSONViewer({ clickableUrls: false });
          viewer.render(this.queryResult);
        } catch (error) {
          console.log('Wrong json format.', error)
        }
      })

    },
  },
}

class JSONViewer {
  constructor(options = {}) {
    this.options = Object.assign({ rootCollapsable: true, clickableUrls: true, bigNumbers: false }, options);
  }

  isCollapsable(arg) {
    return arg instanceof Object && Object.keys(arg).length > 0;
  }

  isUrl(string) {
    const protocols = ['http', 'https', 'ftp', 'ftps'];
    return protocols.some(protocol => string.startsWith(protocol + '://'));
  }

  htmlEscape(s) { return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/'/g, '&apos;').replace(/"/g, '&quot;'); }

  json2html(json) {
    let html = '';
    if (typeof json === 'string') {
      json = this.htmlEscape(json);
      if (this.options.clickableUrls && this.isUrl(json)) {
        html += `<a href="${json}" class="json-string" target="_blank">${json}</a>`;
      } else {
        json = json.replace(/&quot;/g, '\\&quot;');
        html += `<span class="json-string">"${json}"</span>`;
      }
    } else if (typeof json === 'number' || typeof json === 'bigint') {
      html += `<span class="json-literal">${json}</span>`;
    } else if (typeof json === 'boolean') {
      html += `<span class="json-literal">${json}</span>`;
    } else if (json === null) {
      html += '<span class="json-literal">null</span>';
    } else if (Array.isArray(json)) {
      if (json.length > 0) {
        html += '[<ol class="json-array">';
        for (let i = 0; i < json.length; ++i) {
          html += '<li>';
          if (this.isCollapsable(json[i])) {
            html += '<a class="json-toggle"></a>';
          }
          html += this.json2html(json[i]);
          if (i < json.length - 1) {
            html += ',';
          }
          html += '</li>';
        }
        html += '</ol>]';
      } else {
        html += '[]';
      }
    } else if (typeof json === 'object') {
      if (this.options.bigNumbers && (typeof json.toExponential === 'function' || json.isLosslessNumber)) {
        html += `<span class="json-literal">${json.toString()}</span>`;
      } else {
        const keyCount = Object.keys(json).length;
        if (keyCount > 0) {
          html += '{<ul class="json-dict">';
          let count = 0;
          for (const key in json) {
            if (Object.prototype.hasOwnProperty.call(json, key)) {
              const jsonElement = json[key];
              const escapedKey = this.htmlEscape(key);
              const keyRepr = `<span class="json-string-key">"${escapedKey}"</span>`;
              html += '<li>';
              if (this.isCollapsable(jsonElement)) {
                html += `<a class="json-toggle">${keyRepr}</a>`;
              } else {
                html += keyRepr;
              }
              html += ': ' + this.json2html(jsonElement);
              if (++count < keyCount) {
                html += ',';
              }
              html += '</li>';
            }
          }
          html += '</ul>}';
        } else {
          html += '{}';
        }
      }
    }
    return html;
  }

  render(jsonData) {
    const html = this.json2html(jsonData);
    const rootElement = document.getElementById('json-renderer');

    if (this.options.rootCollapsable && this.isCollapsable(jsonData)) {
      rootElement.innerHTML = '<a class="json-toggle"></a>' + html;
    } else {
      rootElement.innerHTML = html;
    }
    const toggleChildren = (element, collapse) => {
      const childToggles = element.querySelectorAll('.json-toggle');
      childToggles.forEach(toggle => {
        const container = toggle.nextElementSibling;
        if (container && (container.classList.contains('json-dict') || container.classList.contains('json-array'))) {
          const isCollapsed = toggle.classList.contains('collapsed');
          if (collapse !== isCollapsed) {
            toggle.click();
          }
        }
      });
    };
    document.addEventListener('click', (e) => {
      if (e.target.matches('.json-toggle')) {
        e.preventDefault();
        const target = e.target;
        const listItem = target.closest('li') || target.parentElement;
        target.classList.toggle('collapsed');
        const container = target.nextElementSibling;
        if (container && (container.classList.contains('json-dict') || container.classList.contains('json-array'))) {
          container.classList.toggle('hidden');
          if (e.ctrlKey) {
            toggleChildren(container, target.classList.contains('collapsed'));
          }
          const siblings = container.parentNode.children;
          for (let i = Array.from(siblings).indexOf(container) + 1; i < siblings.length; i++) {
            if (siblings[i].classList.contains('json-placeholder')) {
              siblings[i].remove();
              i--;
            }
          }
          if (container.classList.contains('hidden')) {
            const count = container.children.length;
            let placeholder = document.createElement('a');
            placeholder.className = 'json-placeholder';
            placeholder.textContent = `${count} ${count > 1 ? 'items' : 'item'}`;
            if (!container.nextElementSibling?.classList.contains('json-placeholder')) {
              container.parentNode.insertBefore(placeholder, container.nextSibling);
            }
          }
        }
      } else if (e.target.matches('.json-placeholder')) {
        e.preventDefault();
        const toggle = e.target.previousElementSibling.previousElementSibling;
        if (toggle) {
          toggle.click();
        }
      }
    });
  }
}
</script>
<style>
.theme--light.v-application {
  background: none !important;
}

.json-highlight {
  padding: 20px;
  /* border-radius: 5px; */
  /* max-height: 700px; */
  max-height: 100vh - 60px;
  /* overflow-y: hidden; */
  /* scrollbar-width: none; */
  /* margin: auto; */
}

.centerArea {
  background-color: rgba(73, 73, 73, 0.5);
  height: 100%;
  width: 100%;
  margin: auto;
  display: flexbox;
}

.queryArea {
  height: 60px;
  margin: 20px;
  color: rgb(206, 206, 206);
}

.showArea {
  height: calc(100vh - 90px);
  color: rgba(255, 255, 255, 0.8);
  font-weight: 700;
  overflow: auto;
}

.showArea::-webkit-scrollbar {
  width: 0;
  /* 隐藏垂直滚动条 */
  height: 0;
  /* 隐藏水平滚动条 */
  display: none;
  /* 完全隐藏滚动条 */
}

.showKeys {
  float: left;
  /* background-color: rgba(37, 37, 37, 0.8); */
  background-color: rgba(27, 27, 27, 0.8);
  height: 100%;
  width: 250px;
  overflow: auto;
}

.showKeys::-webkit-scrollbar,
.showRes::-webkit-scrollbar {
  width: 5px;
  /* 隐藏垂直滚动条 */
  height: 10px;
  /* 隐藏水平滚动条 */
  /* display: none; */
  /* 完全隐藏滚动条 */
}

.showKeys::-webkit-scrollbar-track,
.showRes::-webkit-scrollbar-track {
  background: #3a3a3a;
  /* 轨道背景色 */
  border-radius: 5px;
  /* 圆角 */
}

.showKeys::-webkit-scrollbar-thumb,
.showRes::-webkit-scrollbar-thumb {
  background: #bdccf8;
  /* 滑块颜色 */
  border-radius: 5px;
  /* 圆角 */
}

.showRes {
  background-color: rgba(27, 27, 27, 0.8);
  float: left;
  height: 100%;
  margin-left: 1px;
  width: calc(100% - 251px) !important;
  overflow: auto;
}

.showArea .v-list {
  background: none !important;
  color: rgba(145, 166, 184, 0.8) !important;
}

.v-list-item:hover{
  background-color: rgba(100, 100, 100, 0.8) !important;
}

/****************************JSON-VIEW 样式**************************************/
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  overflow-wrap: break-word;
  display: block;
}

#json-renderer {
  padding: 1em 2em;
}

ul.json-dict,
ol.json-array {
  list-style-type: none;
  margin: 0 0 0 1px;
  border-left: 1px dotted #666;
  padding-left: 2em;
}

.json-string-key {
  color: #7bdcfe;
}

.json-string {
  color: #ce9178;
}

.json-literal {
  color: #b5cea8;
  font-weight: bold;
}

a.json-toggle {
  position: relative;
  color: inherit;
  text-decoration: none;
  cursor: pointer;
}

a.json-toggle:focus {
  outline: none;
}

a.json-toggle:before {
  font-size: 0.8em;
  color: #666;
  content: "\25BC";
  position: absolute;
  display: inline-block;
  width: 1em;
  text-align: center;
  line-height: 1em;
  left: -1.3em;
  top: 1px;
}

a.json-toggle:hover:before {
  color: #aaa;
}

a.json-toggle.collapsed:before {
  content: "\25B6";
}

a.json-placeholder {
  color: #aaa;
  padding: 0 1em;
  text-decoration: none;
  cursor: pointer;
}

a.json-placeholder:hover {
  text-decoration: underline;
}

.hidden {
  display: none;
}</style>