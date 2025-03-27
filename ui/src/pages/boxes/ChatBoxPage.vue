<script setup lang="ts">

import { onMounted, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';
import EventSourceStream from '@server-sent-stream/web';
import { KellerBotChatEvent } from 'src/models/keller_bot_event';

const messages = ref<KellerBotChatEvent[]>([]);

function getStream() {
  Keller_bot.getChatStream().then(async (response) => {
    const stream = response.data;

    const decoder = new EventSourceStream();
    stream.pipeThrough(decoder);

    const reader = decoder.readable.getReader();

    while (true) {
      const { done, value } = await reader.read();

      const chatMessage = KellerBotChatEvent.fromApi(JSON.parse(value!.data));

      if (done) break;
      console.log(chatMessage);
      messages.value.push(chatMessage);
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
  <q-layout>
    <q-page-container>
      <q-page class="bg-black text-white text-h4">
        <div v-for="(message, index) in messages" :key="index">
          <span class="text-bold">{{message.User}}</span> {{message.Message}}
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<style scoped>
</style>
