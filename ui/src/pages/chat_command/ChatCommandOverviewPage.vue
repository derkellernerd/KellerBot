<script setup lang="ts">
import {
  type ApiKellerBotChatCommand,
  KellerBotChatCommand,
} from 'src/models/keller_bot_chat_command';
import { onMounted, ref } from 'vue';
import { emptyPagination } from 'src/helper/common';
import Keller_bot from 'src/service/keller_bot';
import { KellerBotAction } from 'src/models/keller_bot_action';

const chatCommands = ref<KellerBotChatCommand[]>([]);
const selectedChatCommand = ref<KellerBotChatCommand>();
const actions = ref<KellerBotAction[]>([]);
const showDialog = ref(false);

const columns = [
  { name: 'id', label: 'ID', field: 'ID', sortable: true },
  { name: 'command', label: 'Command', field: 'Command', sortable: true },
  { name: 'lastUsed', label: 'Last Used', field: 'LastUsed', sortable: true },
  { name: 'used', label: 'Used', field: 'Used', sortable: true },
  { name: 'timeoutInSeconds', label: 'TimeoutInSeconds', field: 'TimeoutInSeconds', sortable: true },
  { name: 'action', label: 'Action', field: 'Action', sortable: true },
];

function loadChatCommands() {
  Keller_bot.getChatCommands()
    .then((result) => {
      if (result.status === 200) {
        chatCommands.value = result.data.Data.map((s) => KellerBotChatCommand.fromApi(s));
      }
    })
    .catch(() => {
      //TODO: catch
    });
}

function loadActions() {
  Keller_bot.getActions()
    .then((result) => {
      if (result.status === 200) {
        actions.value = result.data.Data.map((s) => KellerBotAction.fromApi(s));
      }
    })
    .catch(() => {
      //TODO: catch
    });
}

function save() {
  if (!selectedChatCommand.value) return;
  if (!selectedChatCommand.value.ID || selectedChatCommand.value.ID == 0) {
    Keller_bot.createChatCommand(selectedChatCommand.value).then((result) => {
      if (result.status === 201) {
        loadChatCommands();
      }
    }).catch(() => {
      //TODO: catch
    });
  }
}

function newChatCommand() {
  selectedChatCommand.value = KellerBotChatCommand.fromApi({} as ApiKellerBotChatCommand);
  showDialog.value = true;
}

onMounted(() => {
  loadActions();
  loadChatCommands();
});
</script>

<template>
  <q-page padding>
    <q-btn color="positive" class="full-width" icon="add" label="add chat command" @click="newChatCommand" />
    <q-table class="q-mt-lg" :rows="chatCommands" :columns="columns" :pagination="emptyPagination" />
    <q-dialog v-model="showDialog" v-if="selectedChatCommand">
      <q-card>
        <q-card-section class="text-h6 text-white bg-primary"> Create Chat Command</q-card-section>
        <q-form @submit="save">
          <q-card-section>
            <q-input label="Command" v-model="selectedChatCommand.Command" />
            <q-select
              label="Action"
              v-model="selectedChatCommand.Action"
              :options="actions"
              map-options
              emit-value
              option-value="ActionName"
              option-label="ActionName"
            />
            <q-input label="Timeout (seconds)" v-model.number="selectedChatCommand.TimeoutInSeconds"/>
          </q-card-section>
          <q-card-section>
            <q-card-actions>
              <q-btn label="cancel" color="negative" v-close-popup type="reset" />
              <q-btn label="save" color="positive" v-close-popup type="submit" />
            </q-card-actions>
          </q-card-section>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style scoped></style>
