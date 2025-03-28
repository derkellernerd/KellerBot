<script setup lang="ts">
import { onMounted, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';
import EventSourceStream from '@server-sent-stream/web';
import { KellerBotAlert, KellerBotAlertType } from 'src/models/keller_bot_alert';

const videoSource = ref<string>();
const gifSource = ref<string>();


function playSound(alert: KellerBotAlert) {
  console.log('playing sound');
  const audio = new Audio(`http://localhost:8080/alert/${alert.ID}`)
  audio.play().then(() => {
    console.log('played');
  }).catch(() => {
    //TODO: fill
  })
}

function playGif(alert: KellerBotAlert) {
  gifSource.value = `http://localhost:8080/alert/${alert.ID}`
  setTimeout(function() {
    gifSource.value = undefined;
  }, 2000)
}

function incomingAlert(alert: KellerBotAlert) {
  switch (alert.Type) {
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

      const alert = KellerBotAlert.fromApi(JSON.parse(value!.data));

      if (done) break;
      incomingAlert(alert);
    }
  } ).catch(() => {
    //TODO: fill
  })
}

onMounted(() => {
  getStream()
})
</script>
<template>
  <img v-if="gifSource" :src="gifSource" alt="gif"/>
  <video v-if="videoSource" :src="videoSource"/>
</template>

<style scoped></style>
