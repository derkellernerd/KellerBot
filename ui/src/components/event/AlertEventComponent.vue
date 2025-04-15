<script setup lang="ts">
import { onMounted, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';
import EventSourceStream from '@server-sent-stream/web';
import { KellerBotAction, KellerBotActionType } from 'src/models/keller_bot_action';
import type { ApiKellerBotAction } from 'src/models/keller_bot_action';

const videoSource = ref<string>();
const gifSource = ref<string>();
const audioSource = ref<string>();
const textAlert = ref<string>();

function playSound(action: KellerBotAction) {

  audioSource.value = `http://localhost:8080/action/${action.ID}`;
  setTimeout(function () {
    audioSource.value = undefined;
  }, action.duration);
}

function playGif(action: KellerBotAction) {
  gifSource.value = `http://localhost:8080/action/${action.ID}`;
  setTimeout(function () {
    gifSource.value = undefined;
  }, action.duration);
}

function showText(action: KellerBotAction) {
  textAlert.value = action.actionText.Text;
  setTimeout(function () {
    textAlert.value = undefined;
  }, action.duration);
}

function incomingAlert(action: KellerBotAction) {
  console.log(`processing alert: ${action.ActionName} -> ${action.ActionType}`);

  switch (action.ActionType) {
    case KellerBotActionType.Text:
      showText(action);
      break;
    case KellerBotActionType.Sound:
      playSound(action);
      break;
    case KellerBotActionType.Gif:
      playGif(action);
      break;
  }
}

function getStream() {
  console.log('start listening');
  Keller_bot.getAlertStream()
    .then(async (response) => {
      const stream = response.data;

      const decoder = new EventSourceStream();
      stream.pipeThrough(decoder);

      const reader = decoder.readable.getReader();

      while (true) {
        const { done, value } = await reader.read();

        console.log('incoming data: ',  value!.data)
        const alert = KellerBotAction.fromApi(JSON.parse(value!.data) as ApiKellerBotAction);

        if (done) break;
        incomingAlert(alert);
      }
      console.log('finish');
    })
    .catch((error) => {
      console.log('Error: ', error);
    })
    .finally(() => {
      getStream();
    });
}

onMounted(() => {
  getStream();
});
</script>
<template>
  <img v-if="gifSource" class="full-height full-width" fit="fill" :src="gifSource" alt="gif" />
  <video v-if="videoSource" :src="videoSource" />
  <audio v-if="audioSource" :src="audioSource" autoplay></audio>
  <div v-if="textAlert" class="text-h2 text-white">
    {{ textAlert }}
  </div>
</template>

<style scoped></style>
