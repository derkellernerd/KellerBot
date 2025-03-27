<script setup lang="ts">
import type { ApiKellerBotChatEvent } from 'src/models/keller_bot_event';
import { ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';

const message = ref<ApiKellerBotChatEvent>({} as ApiKellerBotChatEvent);
const messageInput = ref();

function sendMessage() {
  Keller_bot.createChatEventTest(message.value)
    .then((result) => {
      if (result.status === 204) {
        message.value.Message = '';
        messageInput.value.focus();
      }
    })
    .catch(() => {
      //TODO: fill
    });
}
</script>

<template>
    <q-form @submit="sendMessage" class="row">
      <q-input class="col-2 q-px-md" label="Username" v-model="message.User" />
      <q-input class="col-9 q-px-md" :ref="messageInput" label="Message" v-model="message.Message" />
      <q-btn class="col-1 q-px-md" type="submit" icon="send" label="send" color="primary" />
    </q-form>
</template>

<style scoped></style>
