<script setup lang="ts">
import type { ApiKellerBotActionTypeHttp} from 'src/models/keller_bot_action';
import { KellerBotAction } from 'src/models/keller_bot_action';
import JsonEditorVue from 'json-editor-vue'
import { computed, onMounted } from 'vue';

const action = defineModel({ type: KellerBotAction, required: true });
const actionHttp = computed(() => {
  return action.value.actionHttp;
});

onMounted(() => {
  if (action.value.isNew) {
    action.value.Data = {
      Uri: '',
      HttpMethod: 'GET',
      Payload: {},
    } as ApiKellerBotActionTypeHttp;
  }
});
</script>

<template>
  <q-card class="q-mt-lg" v-if="actionHttp">
    <q-card-section class="text-white text-subtitle2 bg-primary"> Settings </q-card-section>
    <q-card-section>
      <q-input label="URL" v-model="actionHttp.Uri"/>
      <q-select label="Method" v-model="actionHttp.HttpMethod" :options="['GET', 'POST']"/>
      <JsonEditorVue v-if="actionHttp.HttpMethod == 'POST'" v-model="actionHttp.Payload"/>
    </q-card-section>
  </q-card>
</template>

<style scoped>

</style>
