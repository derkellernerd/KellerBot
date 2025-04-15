<script setup lang="ts">
import {
  type ApiKellerBotChatCommand,
  KellerBotChatCommand,
} from 'src/models/keller_bot_chat_command';
import { computed, onMounted, ref } from 'vue';
import { emptyPagination } from 'src/helper/common';
import Keller_bot from 'src/service/keller_bot';
import { KellerBotAction } from 'src/models/keller_bot_action';
import { useQuasar} from 'quasar';
import { isDateSet, timeFormatted } from 'src/helper/functions';

const chatCommands = ref<KellerBotChatCommand[]>([]);
const selectedChatCommand = ref<KellerBotChatCommand>();
const actions = ref<KellerBotAction[]>([]);
const showDialog = ref(false);
const q = useQuasar();
const needle = ref('')

const filteredChatCommands = computed(() => {
  return chatCommands.value.filter(
    (s) => s.Command.toLowerCase().indexOf(needle.value.toLowerCase()) >= 0,
  );
});

const columns = [
  { name: 'command', label: 'Command', field: 'Command', sortable: true },
  {
    name: 'lastUsed',
    label: 'Last Used',
    field: 'LastUsed',
    sortable: true,
    format: (val: Date) => isDateSet(val) ? timeFormatted(val) : 'never',
  },
  { name: 'used', label: 'Used', field: 'Used', sortable: true },
  {
    name: 'timeoutInSeconds',
    label: 'TimeoutInSeconds',
    field: 'TimeoutInSeconds',
    sortable: true,
  },
  { name: 'action', label: 'Action', field: 'Action', sortable: true },
];

function loadChatCommands() {
  Keller_bot.getChatCommands()
    .then((result) => {
      if (result.status === 200) {
        chatCommands.value = result.data.Data.map((s) => KellerBotChatCommand.fromApi(s)).sort(
          (a, b) => a.Command.localeCompare(b.Command),
        );
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
    Keller_bot.createChatCommand(selectedChatCommand.value)
      .then((result) => {
        if (result.status === 201) {
          loadChatCommands();
        }
      })
      .catch(() => {
        //TODO: catch
      });
  }
}

function newChatCommand() {
  selectedChatCommand.value = KellerBotChatCommand.fromApi({} as ApiKellerBotChatCommand);
  showDialog.value = true;
}

function deleteCommand(command: KellerBotChatCommand) {
  q.dialog({
    title: 'Delete Chat Command',
    message: `Do you really want to delete: ${command.Command}?`,
    cancel: true,
    persistent: true,
  }).onOk(() => {
    Keller_bot.deleteChatCommand(command.ID)
      .then((result) => {
        if (result.status === 204) {
          loadChatCommands();
        }
      })
      .catch(() => {
        //TODO: catch
      });
  });
}

onMounted(() => {
  loadActions();
  loadChatCommands();
});
</script>

<template>
  <q-page padding>
    <q-btn
      color="positive"
      class="full-width"
      icon="add"
      label="add chat command"
      @click="newChatCommand"
    />
    <q-input class="q-mt-lg" label="Search" v-model="needle"/>
    <q-table class="q-mt-lg" :rows="filteredChatCommands" :columns="columns" :pagination="emptyPagination">
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
            {{ col.value }}
          </q-td>
          <q-td auto-width>
            <q-btn color="negative" icon="sym_o_delete" @click="deleteCommand(props.row)" />
          </q-td>
        </q-tr>
      </template>
    </q-table>
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
            <q-input
              label="Timeout (seconds)"
              v-model.number="selectedChatCommand.TimeoutInSeconds"
            />
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
