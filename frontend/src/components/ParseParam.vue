<template>
    <v-row>
        <v-col cols="5">
            <v-autocomplete v-model="parseModelVisible" :items="parseModeList" item-title="text" item-value="value"
                density="compact" @update:modelValue="parseFunc"></v-autocomplete>
        </v-col>
        <v-col cols="5">
            <v-autocomplete v-if="parseModel == 1" v-model="queryMsgVisible" :items="msgNameList" density="compact"
                label="消息名称(为空则返回源数据)" @update:modelValue="parseFunc"></v-autocomplete>
        </v-col>
        <!-- <v-col cols="1">
            <v-btn depressed small color="primary" @click="parseFunc">
                解析
            </v-btn>
        </v-col> -->
    </v-row>
</template>
<script>
import { Greet } from '../../wailsjs/go/main/App'
export default {
    name: 'ParseParam',
    props: {
        parseFunc: {
            type: Function,
        },
        parseModel: {
            type: Number,
            default: 2,
        },
        queryMsg: {
            type: String,
            default: '',
        }
    },
    computed: {
        parseModelVisible: {
            get() {
                return this.parseModel;
            },
            set(v) {
                this.$emit('update:parseModel', v);
            },
        },
        queryMsgVisible: {
            get() {
                return this.queryMsg;
            },
            set(v) {
                this.$emit('update:queryMsg', v);
            },
        },
    },
    data() {
        return {
            parseModeList: [
                { text: 'ProtoBuf', value: 1 },
                { text: '源数据', value: 2 },
                { text: 'Msgpack', value: 3 },
            ],
            msgNameList: [],
        }
    },
    mounted() {
        this.getAllMsg()
        this.parseFunc()
    },
    methods: {
        async getAllMsg() {
            Greet('HandleAllMsgName', JSON.stringify({})).then(result => {
                console.log(result)
                this.msgNameList = result.data
            })
        },
    },
}
</script>