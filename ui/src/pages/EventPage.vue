<script setup lang="ts">
import { ref } from 'vue';
import ChatEventComponent from 'components/event/ChatEventComponent.vue';
import ChatEventTestForm from 'components/forms/ChatEventTestForm.vue';
import { emptyPagination } from 'src/helper/common';

const tab = ref('chat');
</script>

<template>
  <q-tabs
    v-model="tab"
    dense
    class="text-grey"
    active-color="primary"
    indicator-color="primary"
    align="justify"
    narrow-indicator
  >
    <q-tab name="chat" label="chat" />
  </q-tabs>

  <q-separator />

  <q-tab-panels v-model="tab" animated>
    <q-tab-panel name="chat">
      <q-card class="q-mb-lg">
        <q-card-section>
          <ChatEventTestForm/>
        </q-card-section>
      </q-card>
      <ChatEventComponent>
        <template v-slot:body="props">
          <q-table :rows="props.messages" :pagination="emptyPagination" hide-pagination />
        </template>
      </ChatEventComponent>
    </q-tab-panel>
  </q-tab-panels>
</template>

<style scoped></style>
