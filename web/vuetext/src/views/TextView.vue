<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'

import axios from "axios"

import { processTextContent } from '@/const/TextUtils.js'
import { gurls } from '@/const/urls.js'

const route = useRoute()

onMounted(() => {
    textId = route.params.id
    console.log("textId=" + textId)

    getText()
})

let textId = "43"
const grespOk = 200

const textFull = reactive({
    id:0,
    type:'',
    content:'',
    createTime:'',
    updateTime:'',
    comments:[],
    tags:[],
})

function getText() {
    let url = gurls.text.get
    url = url.replace(":id", textId)
    
    axios
    .get(url)
    .then(function(response){
        let result = response.data
        if (result.code != grespOk) {
            console.error("Get text failed:" + JSON.stringify(result))
            ElMessage.warning('Get text failed: ' + result.message)
            return
        }

        copyToTextFull(result.data)
    })
    .catch(function(error) {
        console.log("Get text error:")
        console.error(error)
        ElMessage.error('error: ' + error.message)
    })
}

function copyToTextFull(respData) {
    Object.assign(textFull, respData)
}

</script>

<template>
<div class="text-main">
    <span>ID: {{ textFull.id }}</span>
    <span>TYPE: {{ textFull.type }}</span>
    <div class="text-main-content" v-html="processTextContent(textFull.content)"></div>
    <span>{{ textFull.createTime }}</span>
</div>
<div class="text-tag">

</div>
<div class="text-comment">

</div>
</template>

<style scoped>
.text-main, .text-tag  {
    border-bottom: 1px solid black;
}
.text-main-content {
    color: chocolate;
}
</style>