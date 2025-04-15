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
  { name: 'actionName', label: 'Action', field: 'ActionName', sortable: true },
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

function testEvent(event: KellerBotTwitchEvent) {
  Keller_bot.testTwitchEvent(event.ID).then(result => {
    if (result.status === 204) {
      console.log('getestet')
    }
  }).catch(() => {
    //TODO: catch
  })
}

onMounted(() => {
  loadTwitchEvents();
})
</script>

<template>
  <q-page padding>
    <q-btn class="full-width q-mb-lg" color="positive" icon="add" label="add command" @click="showAdd = true"/>
    <q-table :rows="twitchEvents" :columns="columns" :pagination="emptyPagination" hide-pagination>
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th
            v-for="col in props.cols"
            :key="col.name"
            :props="props"
          >
            {{ col.label }}
          </q-th>
          <q-th auto-width />
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td
            v-for="col in props.cols"
            :key="col.name"
            :props="props"
          >
            {{ col.value }}
          </q-td>
          <q-td auto-width>
            <q-btn color="secondary" icon="sym_o_experiment"
                   @click="testEvent(props.row)"/>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <TwitchEventCreateDialog v-model="showAdd"/>
  </q-page>
</template>

<style scoped>

</style>
