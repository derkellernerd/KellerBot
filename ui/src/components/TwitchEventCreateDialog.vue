<script setup lang="ts">
import type { ApiKellerBotTwitchEvent } from 'src/models/keller_bot_twitch_event';
import { ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';
import { KellerBotAction } from 'src/models/keller_bot_action';

const show = defineModel({type: Boolean, default: false})
const twitchEvent = ref<ApiKellerBotTwitchEvent>({} as ApiKellerBotTwitchEvent);
const actions = ref<KellerBotAction[]>([]);

function loadAlerts() {
  Keller_bot.getActions()
    .then((result) => {
      if (result.status === 200) {
        actions.value = result.data.Data.map((s) => KellerBotAction.fromApi(s));
      }
    })
    .catch(() => {
      //TODO: fill
    });
}


function createTwitchEvent() {
  if (!twitchEvent.value) return;

  Keller_bot.createTwitchEvent(twitchEvent.value)
    .then((result) => {
      if (result.status === 201) {
        show.value = false;
      }
    })
    .catch(() => {
      //TODO: Dialog
    });
}

</script>

<template>
  <q-dialog v-model="show" persistent @before-show="loadAlerts">
    <q-card>
      <q-card-section class="text-h6 bg-primary text-white"> Create Alert</q-card-section>
      <q-form @submit="createTwitchEvent">
        <q-card-section>
          <q-input v-model="twitchEvent.TwitchEventSubscription" label="Name" />
          <q-select
            v-model="twitchEvent.ActionName"
            label="Alert"
            :options="actions"
            map-options
            emit-value
            option-label="ActionName"
            option-value="ActionName"
          />
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn color="negative" type="reset" label="cancel" v-close-popup />
            <q-btn color="positive" type="submit" label="create" />
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped>

</style>
