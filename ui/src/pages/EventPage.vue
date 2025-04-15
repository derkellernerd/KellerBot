<script setup lang="ts">
import { onMounted, ref } from 'vue';
import ChatEventTestForm from 'components/forms/ChatEventTestForm.vue';
import { emptyPagination } from 'src/helper/common';
import Keller_bot from 'src/service/keller_bot';
import { KellerBotEvent } from 'src/models/keller_bot_event';
import { date } from 'quasar';

const tab = ref('chat');
const events = ref<KellerBotEvent[]>([]);

const eventColumns = [
  { name: 'eventName', label: 'EventName', field: 'EventName', sortable: true },
  { name: 'source', label: 'Source', field: 'sourceIcon', sortable: true },
  { name: 'eventInfo', label: 'Info', field: 'eventInfo', sortable: false },
  {
    name: 'createdAt',
    label: 'CreatedAt',
    field: 'CreatedAt',
    sortable: true,
    format: (val: Date) => `${date.formatDate(val, 'YYYY-MM-DD HH:mm:ss')}`,
  },
];

function loadEvents() {
  Keller_bot.getEvents()
    .then((result) => {
      if (result.status === 200) {
        events.value = result.data.Data.map((s) => KellerBotEvent.fromApi(s));
      }
    })
    .catch(() => {
      //TODO: catch
    });
}

function replayEvent(event: KellerBotEvent) {
  Keller_bot.replayEvent(event.ID)
    .then((result) => {
      if (result.status === 204) {
        console.log('replayed');
      }
    })
    .catch(() => {
      //TODO: catch
    });
}

onMounted(() => {
  loadEvents();
});
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
    <q-tab name="events" label="Events" />
  </q-tabs>

  <q-separator />

  <q-tab-panels v-model="tab" animated>
    <q-tab-panel name="chat">
      <q-card class="q-mb-lg">
        <q-card-section>
          <ChatEventTestForm />
        </q-card-section>
      </q-card>
    </q-tab-panel>
    <q-tab-panel name="events">
      <q-card class="q-mb-lg">
        <q-card-section>
          <q-btn class="full-width" icon="refresh" @click="loadEvents" />
          <q-table :rows="events" :columns="eventColumns" :pagination="emptyPagination">
            <template v-slot:header="props">
              <q-tr :props="props">
                <q-th v-for="col in props.cols" :key="col.name" :props="props">
                  {{ col.label }}
                </q-th>
                <q-th auto-width />
              </q-tr>
            </template>

            <template v-slot:body="props">
              <q-tr :props="props">
                <q-td v-for="col in props.cols" :key="col.name" :props="props">
                  <template v-if="col.name == 'source'">
                    <q-icon :name="col.value" />
                  </template>
                  <template v-else-if="col.name == 'eventName'">
                    <q-badge>{{ col.value }}</q-badge>
                  </template>
                  <template v-else>
                    {{ col.value }}
                  </template>
                </q-td>
                <q-td auto-width>
                  <q-btn color="secondary" icon="sym_o_replay" @click="replayEvent(props.row)" />
                </q-td>
              </q-tr>
            </template>
          </q-table>
        </q-card-section>
      </q-card>
    </q-tab-panel>
  </q-tab-panels>
</template>

<style scoped></style>
