<script setup lang="ts">
import { OBSWebSocket } from 'obs-websocket-js';
import { onMounted, onUnmounted, ref } from 'vue';

const obs = new OBSWebSocket();
const isConnected = ref(false);
const scenes = ref<string[]>([]);

async function connectToObs() {
  await obs.connect('ws://127.0.0.1:4456');

  isConnected.value = true;
  await loadScenes();
}

function disconnectFromObs() {
  if (!isConnected.value) return;

  void obs.disconnect().then(() => {
    console.log('disconnected');
  });
}

async function loadScenes() {
  if (!isConnected.value) return;
  console.log('getting scenes');
  const result = await obs.call('GetSceneList');

  scenes.value = [];
  result.scenes.forEach((s) => scenes.value.push(s.sceneName as string));
}

async function switchScene(sceneName: string) {
  await obs.call('SetCurrentProgramScene', {sceneName: sceneName});
}

onMounted(async () => {
  await connectToObs();
});

onUnmounted(() => {
  disconnectFromObs();
});
</script>

<template>
  <q-layout>
    <q-page-container>
      <q-page>
        Connected: {{ isConnected }}
        <q-btn icon="refresh" @click="loadScenes" />
        <q-list v-if="isConnected">
          Scenes
          <q-item v-for="scene in scenes" :key="scene">
            <q-item-section><q-btn color="primary" :label="scene" @click="switchScene(scene)"/></q-item-section>
          </q-item>
        </q-list>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<style scoped></style>
