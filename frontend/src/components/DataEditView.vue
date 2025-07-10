<template>
    <v-dialog v-if="needDialog" v-model="editFlagVisible" max-width="600" height="600">
        <v-card>
            <v-card-title>{{ editItemVisible.field }}</v-card-title>
            <v-card-item>
                <ParseParam :parseFunc="parseData" v-model:parseModel="parseModel" v-model:queryMsg="queryMsg"></ParseParam>
            </v-card-item>
            <JsonView :data.sync="editItemVisible.value" ref="jsonView"></JsonView>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn text="关闭" variant="tonal" @click="closeEditView"></v-btn>
                <v-btn v-if="parseModel != 3" color="primary" text="确认" variant="tonal" @click="sureEditItem"></v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
    <div v-else class="notDialog">
        <ParseParam :parseFunc="parseData" v-model:parseModel="parseModel" v-model:queryMsg="queryMsg"></ParseParam>
        <JsonView :data.sync="editItemVisible.value" ref="jsonView"></JsonView>
        <v-spacer></v-spacer>
        <v-btn style="left: 90%; top: 20px;" v-if="parseModel != 3" color="primary" text="修改" @click="sureEditItem"></v-btn>
    </div>
</template>
<script>
import { Greet } from '../../wailsjs/go/main/App'
import JsonView from './JsonView.vue'
import ParseParam from './ParseParam.vue'
export default {
    name: 'DataEditView',
    components: { JsonView, ParseParam },
    props: {
        needDialog: {
            type: Boolean,
            default: true,
        },
        queryKey: {
            type: String,
        },
        dataType: {
            type: String,
        },
        editFlag: {
            type: Boolean,
        },
        editItem: {
            type: Object,
        },
    },
    computed: {
        editFlagVisible: {
            get() {
                return this.editFlag;
            },
            set(v) {
                this.$emit('update:editFlag', v);
            },
        },
        editItemVisible: {
            get() {
                return this.editItem;
            },
            set(v) {
                this.$emit('update:editItem', v);
            },
        },
    },
    mounted() {

    },
    data() {
        return {
            parseModel: 2,
            queryMsg: '',
            oldItem: {},
        }
    },
    mounted() {
        //初始化记录旧数据
        this.oldItem = this.editItem
        console.log("::::::::: ", this.editItem)
    },
    methods: {
        closeEditView() {
            this.editItemVisible = this.oldItem
            this.editFlagVisible = false
        },
        sureEditItem() {
            //TODO 发起修改数据请求
            var params = {
                key: this.queryKey,
                field: this.editItemVisible.field,  //传参,hash:field; list/set:index; zset:member
                oldItem: this.oldItem,              //set传参
                parseMode: this.parseModel + '',
                msg: this.queryMsg,
                // data: JSON.stringify(this.editItemVisible.value),
                data: this.getResText(),
            }
            if (this.parseModel == 1 && params.msg == '') {
                alert("请选择消息类型！")
                return
            }
            console.log("修改数据参数：", params)
            let that = this
            Greet('HandleModifyData', JSON.stringify(params)).then(result => {
                console.log(result)
                if (result.data.error != '') {
                    that.$refs.jsonView.viewer.render(result.data.error)
                    return
                }
                var jsonData = JSON.parse(result.data.data)
                that.editItemVisible.value = jsonData
                that.$refs.jsonView.viewer.render(jsonData)

                this.editFlagVisible = false
            })
        },
        getResText() {
            var resText = this.$refs.jsonView.viewer.getText()
            console.log(resText)
            return JSON.stringify(JSON.parse(resText), null, 0).replace(/\s+/g, '')
            // return resText
        },
        parseData() {
            var params = {
                key: this.queryKey,
                field: this.editItemVisible.field,  //传参,hash:field; list/set:index; zset:member
                oldItem: this.oldItem,              //set传参
                parseMode: this.parseModel + '',
                msg: this.queryMsg,
            }
            console.log(params)
            let that = this
            Greet('HandleParseData', JSON.stringify(params)).then(result => {
                console.log(result)
                if (result.data.error != '') {
                    that.$refs.jsonView.viewer.render(result.data.error)
                    return
                }
                that.editItemVisible.value = result.data.data
                that.$refs.jsonView.viewer.render(result.data.data)
            })
        },
    },
}
</script>

<style>
.notDialog {
    padding-top: 10px;
    /* height: calc(100vh - 250px); */
    height: calc(100vh - 90px);
    width: calc(100vw - 300px);
    background-color: rgba(203, 205, 207, 0.8);
    color: rgb(48, 48, 48);
}
</style>