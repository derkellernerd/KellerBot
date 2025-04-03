<script setup lang="ts">
import type {
  ApiKellerBotAlertCreateRequest,
  ApiKellerBotAlertTypeCompositionCreateRequest,
} from 'src/models/keller_bot_alert';
import { KellerBotAlert, KellerBotAlertType } from 'src/models/keller_bot_alert';
import { computed, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';

const show = defineModel({ type: Boolean, default: false });
const kellerBotAlert = ref<ApiKellerBotAlertCreateRequest>({} as ApiKellerBotAlertCreateRequest);
const composition = ref<ApiKellerBotAlertTypeCompositionCreateRequest>({ AlertNames: [] });
const alerts = ref<KellerBotAlert[]>([]);
const file = ref<File>();
const emits = defineEmits(['created']);

const types = computed(() => {
  return [KellerBotAlertType.Gif, KellerBotAlertType.Composition, KellerBotAlertType.Sound];
});

function loadAlerts() {
  Keller_bot.getAlerts()
    .then((result) => {
      if (result.status === 200) {
        alerts.value = result.data.Data.map((s) => KellerBotAlert.fromApi(s));
      }
    })
    .catch(() => {
      //TODO: fill
    });
}

function createAlert() {
  if (kellerBotAlert.value.Type == KellerBotAlertType.Composition) {
    kellerBotAlert.value.Data = composition.value;
  } else {
    if (!file.value) return;
  }

  Keller_bot.createAlert(kellerBotAlert.value)
    .then((result) => {
      if (result.status === 201) {
        if (kellerBotAlert.value.Type != KellerBotAlertType.Composition) {
          uploadFile(result.data.Data.ID);
        }
      }
    })
    .catch(() => {
      //TODO: fill
    });
}

function uploadFile(alertId: number) {
  if (!file.value) return;

  Keller_bot.uploadAlertFile(alertId, file.value)
    .then((result) => {
      if (result.status === 200) {
        show.value = false;
        emits('created', result.data.Data);
      }
    })
    .catch(() => {
      //TODO: fill
    });
}
</script>

<template>
  <q-dialog v-model="show" persistent @before-show="loadAlerts">
    <q-card>
      <q-card-section class="text-h6 bg-primary text-white"> Create Alert</q-card-section>
      <q-form @submit="createAlert">
        <q-card-section>
          <q-input v-model="kellerBotAlert.Name" label="Name" />
          <q-select v-model="kellerBotAlert.Type" label="Type" :options="types" />
          <q-file
            v-if="kellerBotAlert.Type != KellerBotAlertType.Composition"
            outlined
            v-model="file"
          >
            <template v-slot:prepend>
              <q-icon name="attach_file" />
            </template>
          </q-file>
          <template
            v-if="kellerBotAlert.Type == KellerBotAlertType.Composition && composition != null"
          >
            <q-select
              :key="index"
              v-model="composition.AlertNames[index]"
              v-for="index in composition?.AlertNames.length"
              :options="alerts"
              map-options
              emit-value
              option-value="Name"
              option-label="Name"
            />
            <q-btn label="Add alert" @click="composition.AlertNames.push('')" />
          </template>
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

<style scoped></style>
