<script setup lang="ts">
import {
  type ApiKellerBotActionTypeText,
  KellerBotAction
} from 'src/models/keller_bot_action';
import { computed, onMounted } from 'vue';

const action = defineModel({ type: KellerBotAction, required: true });
const actionText = computed(() => {
  return action.value.actionText;
});

onMounted(() => {
  if (action.value.isNew) {
    action.value.Data = {
      Text: '',
    } as ApiKellerBotActionTypeText;
  }
});
</script>

<template>
  <q-card class="q-mt-lg" v-if="actionText">
    <q-card-section class="text-white text-subtitle2 bg-primary"> Settings </q-card-section>
    <q-card-section>
      <q-input label="Text" v-model="actionText.Text"/>
    </q-card-section>
  </q-card>
</template>

<style scoped>

</style>
