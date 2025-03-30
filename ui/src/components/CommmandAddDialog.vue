<script setup lang="ts">
import { computed, ref } from 'vue';
import type {
  ApiKellerBotCommandAlertAction,
  ApiKellerBotCommandCreateRequest,
  ApiKellerBotCommandMessageAction,
} from 'src/models/keller_bot_command';
import { KellerBotCommandType } from 'src/models/keller_bot_command';
import Keller_bot from 'src/service/keller_bot';
import { KellerBotAlert } from 'src/models/keller_bot_alert';

const show = defineModel({ type: Boolean, default: false });

const commandCreateRequest = ref<ApiKellerBotCommandCreateRequest>();
const messageAction = ref<ApiKellerBotCommandMessageAction>();
const alertAction = ref<ApiKellerBotCommandAlertAction>();

const alerts = ref<KellerBotAlert[]>([]);

function reset() {
  commandCreateRequest.value = {} as ApiKellerBotCommandCreateRequest;
  messageAction.value = {} as ApiKellerBotCommandMessageAction;
  alertAction.value = {} as ApiKellerBotCommandAlertAction;

  loadAlerts()
}

function loadAlerts() {
  Keller_bot.getAlerts().then(result => {
    if (result.status === 200) {
      alerts.value = result.data.Data.map(s => KellerBotAlert.fromApi(s))
    }
  }).catch(() => {
    //TODO: Fill
  })
}

function createCommand() {
  if (!commandCreateRequest.value) return;
  if (!messageAction.value) return;
  if (!alertAction.value) return;

  switch (commandCreateRequest.value.Type) {
    case KellerBotCommandType.MessageAction:
      commandCreateRequest.value.Data = messageAction.value;
      break;
    case KellerBotCommandType.HttpAction:
      return;;
    case KellerBotCommandType.AlertAction:
      commandCreateRequest.value.Data = alertAction.value;
      break;
  }

  Keller_bot.createCommand(commandCreateRequest.value)
    .then((result) => {
      if (result.status === 201) {
        show.value = false;
      }
    })
    .catch(() => {
      //TODO: Dialog
    });
}

const types = computed(() => {
  return [
    KellerBotCommandType.MessageAction,
    KellerBotCommandType.HttpAction,
    KellerBotCommandType.AlertAction,
  ];
});
</script>

<template>
  <q-dialog v-model="show" @before-show="reset">
    <q-card v-if="commandCreateRequest && messageAction && alertAction">
      <q-card-section class="text-h6 bg-primary text-white"> Create Command</q-card-section>
      <q-form @submit="createCommand">
        <q-card-section>
          <q-input v-model="commandCreateRequest.Command" label="Command" />
          <q-select v-model="commandCreateRequest.Type" label="Type" :options="types" />
          <q-input v-model.number="commandCreateRequest.TimeoutInSeconds" label="Timeout (Seconds)" />

          <q-input
            v-if="commandCreateRequest.Type == KellerBotCommandType.MessageAction"
            v-model="messageAction.Message"
            label="Message"
          />
          <q-select
            v-if="commandCreateRequest.Type == KellerBotCommandType.AlertAction"
            v-model="alertAction.Alert"
            label="Alert"
            :options="alerts"
            map-options
            emit-value
            option-label="Name"
            option-value="Name"
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

<style scoped></style>
