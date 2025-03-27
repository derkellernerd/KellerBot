<script setup lang="ts">

import { onMounted, ref } from 'vue';
import { KellerBotCommand } from 'src/models/keller_bot_command';
import Keller_bot from 'src/service/keller_bot';
import { emptyPagination } from 'src/helper/common';
import CommmandAddDialog from 'components/CommmandAddDialog.vue';

const commands = ref<KellerBotCommand[]>([]);
const showAdd = ref(false);

const columns = [
  { name: 'command', label: 'Command', field: 'Command', sortable: true },
  { name: 'commandType', label: 'Type', field: 'typeLabel', sortable: true },
  { name: 'used', label: 'Used', field: 'Used', sortable: true },
]

function loadCommands() {
  Keller_bot.getCommands().then(result => {
    if (result.status === 200) {
      commands.value = result.data.Data.map(s => KellerBotCommand.fromApi(s))
    }
  }).catch(() => {
    //TODO: Dialog
  })
}

onMounted(() => {
  loadCommands();
})
</script>

<template>
<q-page padding>
  <q-btn class="full-width q-mb-lg" color="positive" icon="add" label="add command" @click="showAdd = true"/>
  <q-table :rows="commands" :columns="columns" :pagination="emptyPagination" hide-pagination/>
  <CommmandAddDialog v-model="showAdd"/>
</q-page>
</template>

<style scoped>

</style>
