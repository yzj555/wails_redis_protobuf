<template>
  <div id="outerLayer">
    <v-app class="app">
      <div class="centerArea">
        <div class="queryArea">
          <v-row>
            <v-col cols="5">
              <v-autocomplete v-model="currServer" :items="config.redisServer" item-title="name" item-value="name"
                density="compact" @update:modelValue="changeServer"></v-autocomplete>
            </v-col>
            <v-col cols="5">
              <v-text-field label="Redis-Key" v-model="keyWords" density="compact" @click:append="searchKey"
                append-icon="mdi-magnify"></v-text-field>
            </v-col>
          </v-row>
        </div>
        <div class="showArea">
          <!-- <pre v-html="highlightJSON(queryResult)" class="json-highlight">
                    </pre> -->
          <!-- {{ queryResult }} -->
          <Resizable class="showKeys">
            <div slot="rr">
              <v-list density="compact">
                <v-list-item v-for="(item, i) in keysShow" :key="i" :value="item" active-color="green"
                  @click="selectKey(item)">
                  <v-list-item-title v-text="item"></v-list-item-title>
                </v-list-item>
              </v-list>
            </div>
          </Resizable>
          <!-- 需要根据数据类型展示源数据 -->
          <div v-if="keyType == 'hash'">
            <!-- 123 -->
            <HashDataView :queryKey="queryKey" :data="queryResult"></HashDataView>
          </div>
          <div v-else-if="keyType == 'list' || keyType == 'set'">
            <ListOrSetDataView :queryKey="queryKey" :data="queryResult" :dataType="keyType"></ListOrSetDataView>
          </div>
          <div v-if="keyType == 'zset'">
            <ZsetDataView :queryKey="queryKey" :data="queryResult"></ZsetDataView>
          </div>
          <div v-if="keyType == 'string'">
            <StringDataView :queryKey="queryKey" :data="queryResult"></StringDataView>
          </div>
          <div v-else>
            暂无{{ keyType }}类型处理
          </div>
        </div>
      </div>
    </v-app>
  </div>
</template>

<script>
import { Greet } from '../../wailsjs/go/main/App'

import Resizable from './Resizable.vue'
import HashDataView from './HashDataView.vue'
import ListOrSetDataView from './ListOrSetDataView.vue'
import ZsetDataView from './ZsetDataView.vue'
import StringDataView from './StringDataView.vue'

export default {
  name: "HelloWorld",
  components: { Resizable, HashDataView, ListOrSetDataView, ZsetDataView, StringDataView },
  data() {
    return {
      config: {},
      currServer: '',
      parseModeList: [
        { text: 'ProtoBuf', value: 1 },
        { text: '源数据', value: 2 },
        { text: 'Msgpack', value: 3 },
      ],
      beforeUnloadTime: 0,
      msgNameList: [],
      keyWords: "*",
      keysResult: [],
      keysShow: [],
      keysShowIndex: 0,
      queryKey: "",
      queryField: "",
      parseMode: 1, //默认pb
      queryMsg: "",
      keyType: "", //当前查询的key的类型
      queryResult: [],
    }
  },
  mounted() {
    this.getConfig()
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
    clearData() {
      this.keyWords = '*'
      this.keysResult = []
      this.keysShow = []
      this.queryKey = ""
      this.queryField = ""
      this.queryMsg = ""
      this.keyType = ""
      this.queryResult = {}
      this.jsonRender()
    },
    async getConfig() {
      Greet('HandleConfig', JSON.stringify({})).then(result => {
        console.log(result)
        this.config = result.data
        this.currServer = this.config.currentRedis
      })
    },
    async changeServer() {
      Greet('HandleChangeServer', JSON.stringify({ name: this.currServer })).then(result => {
        console.log(result)
        this.clearData()
      })
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
        // this.keysResult = result.data

        this.keysResult = Object.freeze(result.data);
        this.keysShow = []
        this.keysShowIndex = 0
        // this.keysResult =  Object.freeze(result.data)
        this.loadChunk()
      })
    },
    loadChunk() {
      const chunk = this.keysResult.slice(this.keysShowIndex, this.keysShowIndex + 1);
      this.keysShow.push(...chunk);
      this.keysShowIndex += 1;
      if (this.keysShowIndex < this.keysResult.length) {
        requestAnimationFrame(this.loadChunk)
      }
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
        this.keyType = result.data.dataType
        this.queryResult = result.data.data
        console.log(result)
      })

    },
  },
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

thead {
  color: rgb(82, 82, 102) !important;
  font-size: larger;
  font-weight: 800;
}

.v-table {
  background-color: rgba(203, 205, 207, 0.8) !important;
  color: #1b2a36 !important;
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
  color: rgb(255, 255, 255);
}

.showArea {
  height: calc(100vh - 90px);
  color: rgba(255, 255, 255, 0.8);
  font-weight: 700;
  overflow: auto;
  display: flex;
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
  min-width: 250px;
  overflow: auto;
  overflow-y: auto;
  overflow-x: hidden;
  /* resize: horizontal;
  cursor: ew-resize; */
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

.v-list-item:hover {
  background-color: rgba(100, 100, 100, 0.8) !important;
}

</style>