<script setup lang="ts">
import { onMounted, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';
import EventSourceStream from '@server-sent-stream/web';
import type { ApiKellerBotAlert } from 'src/models/keller_bot_alert';
import { KellerBotAlertType,  KellerBotAlert } from 'src/models/keller_bot_alert';

const videoSource = ref<string>();
const gifSource = ref<string>();
const audioSource = ref<string>();

function playSound(alert: KellerBotAlert) {
  audioSource.value = `http://localhost:8080/alert/${alert.ID}`
  setTimeout(function() {
    audioSource.value = undefined;
  }, 2000)
}

function playGif(alert: KellerBotAlert) {
  gifSource.value = `http://localhost:8080/alert/${alert.ID}`
  setTimeout(function() {
    gifSource.value = undefined;
  }, 2000)
}

function incomingAlert(alert: KellerBotAlert) {
  console.log(`processing alert: ${alert.Name} -> ${alert.Type}`)
  switch (alert.Type) {
    case KellerBotAlertType.Composition:
      alert.composition.Alerts.forEach((childAlert) => {
        incomingAlert(KellerBotAlert.fromApi(childAlert));
      })
      break;
    case KellerBotAlertType.Sound:
      playSound(alert)
      break;
    case KellerBotAlertType.Gif:
      playGif(alert)
      break;
  }
}

function getStream() {
  console.log('start listening')
  Keller_bot.getAlertStream().then(async (response) => {
    const stream = response.data;

    const decoder = new EventSourceStream();
    stream.pipeThrough(decoder);

    const reader = decoder.readable.getReader();

    while (true) {
      const { done, value } = await reader.read();

      const alerts = (JSON.parse(value!.data) as ApiKellerBotAlert[]).map(s => KellerBotAlert.fromApi(s));

      if (done) break;
      alerts.forEach(s => {
        incomingAlert(s);
      })
    }
    console.log('finish')
  } ).catch(() => {
    //TODO: fill
  })
}

onMounted(() => {
  getStream()
})
</script>
<template>
  <img v-if="gifSource" class="full-height full-width" fit="fill" :src="gifSource" alt="gif" />
  <video v-if="videoSource" :src="videoSource" />
  <audio v-if="audioSource" :src="audioSource" autoplay></audio>
</template>

<style scoped></style>


