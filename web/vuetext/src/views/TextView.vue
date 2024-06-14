<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'

import axios from "axios"

import { processTextContent, processComment } from '@/const/TextUtils.js'
import { gurls } from '@/const/urls.js'

const route = useRoute()

onMounted(() => {
    // string to number
    textId = +route.params.id

    getText()
})

let textId = 0

const grespOk = 200
const formLabelWidth = '80px'

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

// add comment

const showAddComment = ref(false)

const addCommentForm = reactive({
    textId: 0,
    comment:'',
})

function toggleAddComment() {
    addCommentForm.textId = textId

    showAddComment.value = !showAddComment.value
}

function onAddComment() {
    addCommentForm.comment = addCommentForm.comment.trim()
    if (Object.is(addCommentForm.comment, '')) {
        ElMessage.warning("Comment is empty.")
        return
    }

    httpAddComment()
}

function httpAddComment() {
    let url = gurls.comment.add
    axios
    .post(url, addCommentForm)
    .then(function(response) {
        let result = response.data
        let rcode = result.code
        if (! Object.is(rcode, grespOk)) {
            ElMessage.warning("Add comment failed: " + JSON.stringify(result))
            return
        }

        ElMessage.success('Add comment success[' + result.data + ']🎈')

        // update
        getText()

        // reset form data
        addCommentForm.comment = ''
    })
    .catch(function(error) {
        console.error('Add comment error:')
        console.error(error)
        ElMessage.error('Add comment error: ' + error.message)
    })
}

// delete comment
function onDeleteComment(cmt) {
    ElMessageBox.confirm("Confirm the deletion of the comment?")
    .then(() => {
        httpDeleteComment(cmt)
    })
    .catch((e) => {
        // e is cancel
    })
}

function httpDeleteComment(cmt) {
    let url = gurls.comment.del
    url = url.replaceAll(":textId", cmt.textId)
    url = url.replaceAll(":id", cmt.id)
    
    axios
    .delete(url)
    .then(function(response) {
        let result = response.data;
        if (result.code != grespOk) {
            ElMessage.error('Delete comment failed: ' + JSON.stringify(result))
            return;
        }

        ElMessage.success('Delete comment success🎈')

        // get Text again
        getText()
    })
    .catch(function(error) {
        console.error("Delete comment error: ")
        console.error(error);
        ElMessage.error('Delete comment error: ' + error.message);
    });
}

</script>

<template>
    <el-link type="primary" href="/">goto Home Page</el-link>
<div class="text-main">
    <span>ID: {{ textFull.id }}</span>
    <span>TYPE: {{ textFull.type }}</span>
    <div class="text-main-content" v-html="processTextContent(textFull.content)"></div>
    <span>{{ textFull.createTime }}</span>
</div>
<div class="text-tag">
    <div style="font-weight: bold;">Tags:{{ textFull.tags.length }}</div>
    <el-link type="success" v-for="(item,index) in textFull.tags" class="list-tag-name">{{ item.name }}</el-link>
</div>
<div class="text-comment">
    <div style="font-weight: bold;">Comments:{{ textFull.comments.length }}</div>
    <el-button type="primary" @click="toggleAddComment">Add Comment</el-button>
    <div v-show="showAddComment" style="border: 1px solid #00FFFF;">
        <el-form>
            <el-form-item label="Text#ID" :label-width="formLabelWidth">
                <el-text>{{ addCommentForm.textId }}</el-text>
            </el-form-item>
            <el-form-item label="Comment" :label-width="formLabelWidth">
                <el-input v-model="addCommentForm.comment"
                    type="textarea"
                    placeholder="Input comment here"
                    autocomplete="off"
                    :rows="5"
                    maxlength="200"
                    show-word-limit
                    @keyup.ctrl.enter="onAddComment"
                    resize="none"
                    />
            </el-form-item>
            <el-form-item :label-width="formLabelWidth">
                <el-button type="success" @click="onAddComment">Submit</el-button>
            </el-form-item>
        </el-form>
    </div>
    <div v-for="(item,index) in textFull.comments">
        <span>{{ item.createTime }}</span>
        <div v-html="processComment(item.comment)" class="list-comment-item"></div>
        <div>
            <Delete @click="onDeleteComment(item)" class="icon-delete-tag" />
        </div>
    </div>
</div>
</template>

<style scoped>
.text-main, .text-tag  {
    border-bottom: 1px solid black;
}
.text-main-content {
    color: chocolate;
}

.list-tag-name {
    display: inline-block;
    width: 10em;
    text-align: left;
}

.list-comment-item {
    color: chocolate;
}

.icon-delete-tag {
    color: red;
    width: 1em;
    height: 1em; 
    margin-right: 8px;
}

</style>