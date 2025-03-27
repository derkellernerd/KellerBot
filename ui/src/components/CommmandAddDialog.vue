<script setup lang="ts">

import { ref } from 'vue';
import type {
  ApiKellerBotCommandCreateRequest,
  ApiKellerBotCommandMessageAction,
} from 'src/models/keller_bot_command';
import {KellerBotCommand} from 'src/models/keller_bot_command';
import Keller_bot from 'src/service/keller_bot';

const show = defineModel({type: Boolean, default: false})

const commandCreateRequest = ref<ApiKellerBotCommandCreateRequest>()
const messageAction = ref<ApiKellerBotCommandMessageAction>()

function reset() {
  commandCreateRequest.value = {} as ApiKellerBotCommandCreateRequest
  messageAction.value = {} as ApiKellerBotCommandMessageAction
}

function createCommand() {
  if (!commandCreateRequest.value) return
  if (!messageAction.value) return

  commandCreateRequest.value.Data = messageAction.value;

  Keller_bot.createCommand(commandCreateRequest.value).then((result) => {
    if (result.status === 201) {
      show.value = false;
    }
  }).catch(() => {
    //TODO: Dialog
  })
}

</script>

<template>
  <q-dialog v-model="show" @before-show="reset">
    <q-card v-if="commandCreateRequest && messageAction">
      <q-card-section class="text-h6 bg-primary text-white"> Create Command </q-card-section>
      <q-form @submit="createCommand">
        <q-card-section>
          <q-input v-model="commandCreateRequest.Command" label="Command"/>
          <q-select v-model="commandCreateRequest.Type" label="Type" :options="KellerBotCommand.commandTypes"/>
          <q-input v-model="messageAction.Message" label="Message"/>
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn color="negative" type="reset" label="cancel" v-close-popup/>
            <q-btn color="positive" type="submit" label="create"/>
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
