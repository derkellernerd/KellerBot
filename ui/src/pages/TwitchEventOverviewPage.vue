<script setup lang="ts">

import { onMounted, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';
import { emptyPagination } from 'src/helper/common';
import { KellerBotTwitchEvent } from 'src/models/keller_bot_twitch_event';
import TwitchEventCreateDialog from 'components/TwitchEventCreateDialog.vue';

const twitchEvents = ref<KellerBotTwitchEvent[]>([]);
const showAdd = ref(false);

const columns = [
  { name: 'twitchEventSubscription', label: 'TwitchEventSubscription', field: 'TwitchEventSubscription', sortable: true },
  { name: 'alertName', label: 'Alert', field: 'AlertName', sortable: true },
]

function loadTwitchEvents() {
  Keller_bot.getTwitchEvents().then(result => {
    if (result.status === 200) {
      twitchEvents.value = result.data.Data.map(s => KellerBotTwitchEvent.fromApi(s))
    }
  }).catch(() => {
    //TODO: Dialog
  })
}

onMounted(() => {
  loadTwitchEvents();
})
</script>

<template>
  <q-page padding>
    <q-btn class="full-width q-mb-lg" color="positive" icon="add" label="add command" @click="showAdd = true"/>
    <q-table :rows="twitchEvents" :columns="columns" :pagination="emptyPagination" hide-pagination/>
    <TwitchEventCreateDialog v-model="showAdd"/>
  </q-page>
</template>

<style scoped>

</style>
