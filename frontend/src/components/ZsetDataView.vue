<template>
    <div>
        <v-data-table :headers="headers" :items="data" style="height: calc(100vh - 90px); width: calc(100vw - 250px);">
            <template v-slot:item.action="{ item }">
                <div class="d-flex ga-2 justify-center">
                    <v-icon color="medium-emphasis" icon="mdi-pencil" size="small" @click="showEditDialog(item)"></v-icon>
                </div>
            </template>
        </v-data-table>
        <DataEditView :queryKey="queryKey" dataType="zset" v-model:editFlag="editFlag" v-model:editItem="editItem">
        </DataEditView>
    </div>
</template>
<script>
import DataEditView from './DataEditView.vue'
export default {
    name: 'ZsetDataView',
    components: { DataEditView },
    props: {
        queryKey: {
            type: String,
        },
        data: {
            type: Array,
            default: [],
        }
    },
    data() {
        return {
            headers: [
                { title: 'member', align: 'start', key: 'field' },
                { title: 'score', align: 'start', key: 'value', maxWidth: '550', nowrap: true },
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
<style></style>