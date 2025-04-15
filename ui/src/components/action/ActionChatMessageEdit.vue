<script setup lang="ts">
import {
  type ApiKellerBotActionTypeChatMessage,
  KellerBotAction
} from 'src/models/keller_bot_action';
import { computed, onMounted } from 'vue';

const action = defineModel({ type: KellerBotAction, required: true });
const actionChatMessage = computed(() => {
  return action.value.actionChatMessage;
});

onMounted(() => {
  if (action.value.isNew) {
    action.value.Data = {
      ChatMessage: '',
    } as ApiKellerBotActionTypeChatMessage;
  }
});
</script>

<template>
  <q-card class="q-mt-lg" v-if="actionChatMessage">
    <q-card-section class="text-white text-subtitle2 bg-primary"> Settings </q-card-section>
    <q-card-section>
      <q-input label="Text" v-model="actionChatMessage.ChatMessage"/>
    </q-card-section>
  </q-card>
</template>

<style scoped>

</style>
