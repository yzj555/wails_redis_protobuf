<template>
    <div>
        <v-data-table :headers="headers" :items="data" style="height: calc(100vh - 90px); width: calc(100vw - 250px);">
            <template v-slot:item.action="{ item }">
                <div class="d-flex ga-2 justify-center">
                    <v-icon color="medium-emphasis" icon="mdi-pencil" size="small" @click="showEditDialog(item)"></v-icon>
                </div>
            </template>
        </v-data-table>

        <DataEditView :queryKey="queryKey" dataType="hash" v-model:editFlag="editFlag" v-model:editItem="editItem">
        </DataEditView>
    </div>
</template>
<script>
// import { Greet } from '../../wailsjs/go/main/App'
// import JsonView from './JsonView.vue'
// import ParseParam from './ParseParam.vue'
import DataEditView from './DataEditView.vue'

export default {
    name: "HashDataView",
    components: { DataEditView },
    props: {
        queryKey: {
            type: String,
        },
        data: {
            type: Array,
            default: []
        }
    },
    data() {
        return {
            headers: [
                { title: 'field', align: 'start', key: 'field' },
                { title: 'value', align: 'start', key: 'value', maxWidth: '350', nowrap: true },
                { title: 'action', align: 'center', key: 'action' },
            ],
            editFlag: false,
            editItem: {},
        }
    },
    methods: {
        showEditDialog(item) {
            this.editFlag = true
            this.editItem = item
        },
    },
}
</script>