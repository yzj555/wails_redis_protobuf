<template>
    <div class="resizable-container" ref="container">
        <div class="content" :style="{ width: width + 'px' }">
            <slot></slot>
        </div>
        <div class="drag-handle" @mousedown="startDrag"></div>
    </div>
</template>
  
<script setup>
import { ref } from 'vue';

const width = ref(300);
const container = ref(null);
let startX = 0;

const startDrag = (e) => {
    startX = e.clientX;
    window.addEventListener('mousemove', drag);
    window.addEventListener('mouseup', stopDrag);
};

const drag = (e) => {
    width.value += e.clientX - startX;
    startX = e.clientX;
};

const stopDrag = () => {
    window.removeEventListener('mousemove', drag);
};
</script>
  
<style>
.resizable-container {
    position: relative;
    display: inline-block;
}

.drag-handle {
    position: absolute;
    right: -4px;
    top: 0;
    bottom: 0;
    width: 8px;
    resize: horizontal;
    cursor: col-resize;
    background: rgba(0, 0, 0, 0.1);
    z-index: 99;
}
</style>
  