<script setup lang="ts">
import type { ApiKellerBotAlert } from 'src/models/keller_bot_alert';
import {KellerBotAlertType} from 'src/models/keller_bot_alert';
import { computed, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';

const show = defineModel({ type: Boolean, default: false });
const kellerBotAlert = ref<ApiKellerBotAlert>({} as ApiKellerBotAlert)
const file = ref<File>();
const emits = defineEmits(['created'])

const types = computed(() => {
  return [
    KellerBotAlertType.Gif
  ]
})

function createAlert() {
  if (!file.value) return;

  Keller_bot.createAlert(kellerBotAlert.value).then(result => {
    if (result.status === 201) {
      uploadFile(result.data.Data.ID)
    }
  }).catch(() => {
    //TODO: fill
  })
}

function uploadFile(alertId: number) {
  if (!file.value) return;

  Keller_bot.uploadAlertFile(alertId, file.value).then((result) => {
    if (result.status === 200) {
      show.value = false;
      emits('created', result.data.Data);
    }
  }).catch(() => {
    //TODO: fill
  })
}
</script>

<template>
  <q-dialog v-model="show" persistent>
    <q-card>
      <q-card-section class="text-h6 bg-primary text-white"> Create Alert </q-card-section>
      <q-form @submit="createAlert">
        <q-card-section>
          <q-input v-model="kellerBotAlert.Name" label="Name" />
          <q-select v-model="kellerBotAlert.Type" label="Type" :options="types"/>
          <q-file outlined v-model="file">
            <template v-slot:prepend>
              <q-icon name="attach_file" />
            </template>
          </q-file>
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
