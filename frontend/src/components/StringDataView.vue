<template>
    <div>
        <DataEditView :queryKey="queryKey" dataType="string" :needDialog="false" v-model:editItem="editItem">
        </DataEditView>

    </div>
</template>
<script>
import DataEditView from './DataEditView.vue'
export default {
    name: 'StringDataView',
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
            editItem: {},
        }
    },
    mounted() {
        this.editItem.value = this.data
    },
}
</script>