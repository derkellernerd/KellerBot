<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { KellerBotChatEvent } from 'src/models/keller_bot_event';
import Keller_bot from 'src/service/keller_bot';
import EventSourceStream from '@server-sent-stream/web';

const messages = ref<KellerBotChatEvent[]>([]);

function getStream() {
  console.log('start listening')
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
<slot name="body" :messages="messages"></slot>
</template>

<style scoped>

</style>
