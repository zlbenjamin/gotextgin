<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import axios from "axios"

onMounted(() => {
    pageFind()
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
        del: '/api/text/comment/:id'
    }
}

const grespOk = 200

// Paging

// request data
const pageFindForm = reactive({
    pageNo: 1,
    pageSize: 10,
    kwContent: '',
    type: '',
    tags: []
})

// response data
const pfPageNo = ref(1)
const pfPageSize = ref(10)
const pfTotal = ref(0)
const pfTotalPage = ref(0)
const pageFindList = reactive([])

function pageFind() {
    let url = gurls.text.page
    axios.post(url, pageFindForm)
    .then(function(response){
        let result = response.data
        if (result.code != grespOk) {
            console.error("Paging failed:" + JSON.stringify(result))
            ElMessage.warning('Paging failed: ' + result.message)
            return
        }

        pfPageNo.value = result.data.pageNo
        pfPageSize.value = result.data.pageSize
        pfTotal.value = result.data.total
        pfTotalPage.value = result.data.totalPage

        // // TextView-BIvkEUqG.js:2 TypeError: Assignment to constant variable.
        // // pageFindList = result.data.list
        fillListAfterClear(pageFindList, result.data.list)
    })
    .catch(function(error) {
        console.error(error)
        ElMessage.error('error: ' + error.message)
    })
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

function handleSizeChange(num) {
    pageFindForm.pageSize = num
    needUpdate.value++
}
function handleCurrentChange(num) {
    pageFindForm.pageNo = num
    needUpdate.value++
}

// Add a text
const showAdd = ref(false)

const tags0 = ref('')
const addForm = reactive({
    type:'',
    content: '',
    tags: []
})

function addMsg() {
    // validate
    addForm.content = addForm.content.trim()
    if (Object.is(addForm.content, '')) {
        ElMessage.warning('Content is blank.')
        // focus again
        // todo
        return
    }

    // tags0 to addForm.tags
    fillTagsArray(tags0.value, addForm.tags)

    // send request
    let url = gurls.text.add
    axios.post(url, addForm)
    .then(function(response){
        let result = response.data
        if (result.code != grespOk) {
            ElMessage.warning('Add failed: ' + result.message)
            return
        }

        ElMessage.success('Add success🎈')
        needUpdate.value++

        clearAddForm()
    })
    .catch(function(error) {
        console.error(error)
        ElMessage.error('Add error: ' + error.message)
    })
}

// Clearn add form after add success
function clearAddForm() {
    tags0.value = ''

    addForm.type = ''
    addForm.content = ''
    // addForm.tags
    addForm.tags.splice(0, addForm.tags.length)
}

// more search conditions
const showSearch = ref(false)

const searchTags0 = ref('')
const searchForm = reactive({
    // trimmed
    type: '',
    // not trimmed
    kwContent: '',
    // parsed from searchTags0
    tags: []
})

// Start a search
// replace fields in 
function startSearch() {
    pageFindForm.type = searchForm.type

    pageFindForm.kwContent = searchForm.kwContent
    
    fillTagsArray(searchTags0.value, pageFindForm.tags)

    needUpdate.value++
}

function fillTagsArray(str, targetArr) {
    // clear old values
    if (targetArr.length > 0) {
        targetArr.splice(0, targetArr.length)
    }
    
    // str to targetArr
    if (str.length > 0) {
        // pre: Replace consecutive spaces
        let re = new RegExp(" +", "g")
        str = str.replaceAll(re, ' ')
        str = str.trim()

        let arr = str.split(' ')
        if (arr.length > 5) {
            ElMessage.warning('At most 5 tags.')
            // focus again
            // todo
            return
        }
        if (arr.length > 0) {
            // fill targetArr
            // pre: clear targetArr
            for (let i=0; i<arr.length; i++) {
                let item = arr[i].trim()
                if (item.length > 10) {
                    ElMessage.warning('Max length of a tag: 10 characters.')
                    // focus again
                    // todo
                    return
                }
                // check duplicate tag
                if (targetArr.indexOf(item) >= 0) {
                    ElMessage.warning('Duplicated tag: ' + item)
                    // focus again
                    // todo
                    return
                }
                targetArr.push(item)
            }
        }
    }
}

const needUpdate = ref(0)
watch(needUpdate, (val) => {
    pageFind()
})

// Process content of a text
function processTextContent(content) {
    if (Object.is(content, null) || Object.is(content, undefined)) {
        return ''
    }

    let pcontent = content
    pcontent = htmlEncode(pcontent)

    // replace \n
    pcontent = pcontent.replaceAll("\n", "<br>")

    // replace url in content
    // Note, there are some bugs TODO
    let urlg = new RegExp("http[s]://[a-zA-Z0-9\.\?/#=_:\\-\\+%]+", "ig")
    pcontent = pcontent.replaceAll(urlg, replaceUrl)

    return pcontent
}

// HTML special characters
function htmlEncode(content) {
    content = content.replaceAll("&", "&amp;")
			
    content = content.replaceAll("<", "&lt;")
    content = content.replaceAll(">", "&gt;")
    content = content.replaceAll("\"", "&quot;")
    content = content.replaceAll(" ", "&nbsp;")

    content = content.replaceAll("\t", "&emsp;")

    return content
}

// Replace matched url with <a> element
function replaceUrl(url) {
    return "<a href='" + url + "' target='_blank' class='content-url' >" + decodeURI(url) + "</a>"
}

</script>

<template>
    <div>
        <el-button type="primary" @click="showAdd = !showAdd">Add</el-button>
        <el-button type="success" @click="showSearch = !showSearch">Search</el-button>
    </div>
    <div v-show="showAdd">
        <el-form :model="addForm" label-width="auto" style="max-width: 600px">
            <el-form-item label="">
                <el-input v-model.trim="addForm.type" 
                    placeholder="Text type"
                    minlength="1" maxlength="10" show-word-limit
                    clearable />
            </el-form-item>
            <el-form-item label="">
                <el-input v-model="addForm.content"
                type="textarea"
                placeholder="Input content here"
                autocomplete="off"
                :rows="20"
                maxlength="10000"
                show-word-limit
                @keyup.ctrl.enter="addMsg"
                />
            </el-form-item>
            <el-form-item>
                <el-input v-model="tags0" 
                    placeholder="Tags separated by spaces"
                    maxlength="55" show-word-limit
                    clearable />
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="addMsg">Submit</el-button>
            </el-form-item>
        </el-form>
    </div>
    <div v-show="showSearch">
        <h1>Search text</h1>
        <el-form :model="searchForm" label-width="auto" style="max-width: 600px">
            <el-form-item label="">
                <el-input v-model.trim="searchForm.type" 
                    placeholder="Text type (full name)"
                    minlength="1" maxlength="10" show-word-limit
                    clearable />
            </el-form-item>
            <el-form-item label="">
                <el-input v-model="searchForm.kwContent" 
                    placeholder="Key word of the content"
                    minlength="1" maxlength="20" show-word-limit
                    clearable />
            </el-form-item>
            <el-form-item>
                <el-input v-model="searchTags0" 
                    placeholder="Tags separated by spaces (AND)"
                    maxlength="55" show-word-limit
                    clearable />
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="startSearch">Submit</el-button>
            </el-form-item>
        </el-form>
    </div>
<div style="display: flex; justify-content: center; align-items: center;">
                <el-pagination
                    style=""
                    v-model:current-page="pfPageNo"
                    v-model:page-size="pfPageSize"
                    :page-sizes="[10, 25, 50, 100]"
                    :small="small"
                    :disabled="disabled"
                    :background="background"
                    layout="total, sizes, prev, pager, next, jumper"
                    :total="pfTotal"
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange"
                    />
            </div>
<div v-for="(item,index) in pageFindList" class="content-in-list" :class="{
        'content-in-list-0':(index%2 == 0),
        'content-in-list-1':(index%2 == 1)
    }">
    <span class="tf tf-id">{{ item.id }}</span>
    <div v-html="processTextContent(item.content)" class="tf tf-content"></div>
    <span class="tf tf-type">#{{ item.type }}</span>
    <span class="tf tf-time">{{ item.createTime }}</span>
    <div>
        <span style="font-weight: bold;">Tags: ({{ item.tags.length }})</span>
        <br>
        <el-link class="tf tf-tag" v-for="(tag,idx2) in item.tags" @click="">{{ tag.name }}</el-link>
    </div>
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
.tf-tag {
    width: 120px;
    color: red;
    text-align: center;
}
</style>