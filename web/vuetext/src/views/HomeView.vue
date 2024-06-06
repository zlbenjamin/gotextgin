<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import axios from "axios"

onMounted(() => {
    pageFind()
})


// Paging

const pageFindForm = reactive({
    pageNo: 1,
    pageSize: 10,
    kwContent:'',
    type:''
})

const pfPageNo = ref(1)
const pfPageSize = ref(10)
const pfTotal = ref(0)
const pfTotalPage = ref(0)
const pageFindList = reactive([])

function pageFind() {
    let url = '/api/text/page'
    axios.post(url, pageFindForm)
    .then(function(response){
        let result = response.data;
        console.log(result)
        
        pfPageNo.value = result.data.pageNo;
        pfPageSize.value = result.data.pageSize;
        pfTotal.value = result.data.total;
        pfTotalPage.value = result.data.totalPage;

        // // TextView-BIvkEUqG.js:2 TypeError: Assignment to constant variable.
        // // pageFindList = result.data.list;
        fillListAfterClear(pageFindList, result.data.list)
    })
    .catch(function(error) {
        console.error(error);
        ElMessage.error('error: ' + error.message);
    });
}

function fillListAfterClear(target, source) {
    if (target.length > 0) {
        target.splice(0, target.length)
    }

    if (source.length > 0) {
        let i = 0
        for (; i<source.length; i++) {
            target.push(source[i])
        }
    }
}

</script>

<template>
<div v-for="(item,index) in pageFindList" class="content-in-list" :class="{
        'content-in-list-0':(index%2 == 0),
        'content-in-list-1':(index%2 == 1)
    }">
    <span class="tf tf-id">{{ item.id }}</span>
    <span class="tf tf-content">{{ item.content }}</span>
    <span class="tf tf-type">#{{ item.type }}</span>
    <span class="tf tf-time">{{ item.createTime }}</span>
</div>
</template>

<style scoped>
.content-in-list {
    width: 98%;
    margin: auto;
    margin-bottom: 2px;
    padding: 4px;
    border-radius: 0.5em;
}
.content-in-list-0 {
    background-color: #D8D8D8 ;
}
.content-in-list-1 {
    background-color:#8888CC;
    color: #fff;
}

.tf {
    display: inline-block;
}
.tf-id {
    width: 40px;
}
.tf-content {
    width: 600px;
	word-break: break-all;
	word-wrap: break-word;
}
.tf-type {
    width: 120px;
}
.tf-time {
    width: 180px;
}
</style>