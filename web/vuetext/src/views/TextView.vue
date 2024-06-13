<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'

import axios from "axios"

const route = useRoute()

onMounted(() => {
    textId = route.params.id
    console.log("textId=" + textId)
    
    getText()
})

const gurls = {
    text: {
        add: '/api/text/add',
        get: '/api/text/:id',
        del: '/api/text/:id',
        page: '/api/text/page',
    },
    comment: {
        add: '/api/text/comment/add',
        del: '/api/text/comment/:textId/:id'
    }
}

let textId = "43"
const grespOk = 200

const textFull = reactive({})

function getText() {
    let url = gurls.text.get
    url = url.replace(":id", textId)
    console.log("todo url=" + url)
    axios
    .get(url)
    .then(function(response){
        let result = response.data
        if (result.code != grespOk) {
            console.error("Get text failed:" + JSON.stringify(result))
            ElMessage.warning('Get text failed: ' + result.message)
            return
        }

        console.log("todo text=" + JSON.stringify(result.data))
        textFull=result.data
    })
    .catch(function(error) {
        console.log("Get text error:")
        console.error(error)
        ElMessage.error('error: ' + error.message)
    })
}

</script>

<template>
<div>
    todo for text detail...
</div>
</template>

<style scoped>

</style>