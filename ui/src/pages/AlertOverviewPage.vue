<script setup lang="ts">

import { onMounted, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';
import { emptyPagination } from 'src/helper/common';
import { KellerBotAlert } from 'src/models/keller_bot_alert';

const alerts = ref<KellerBotAlert[]>([]);

const columns = [
  { name: 'id', label: 'ID', field: 'ID', sortable: true },
  { name: 'name', label: 'Name', field: 'Name', sortable: true },
  { name: 'alertType', label: 'Type', field: 'Type', sortable: true },
  { name: 'used', label: 'Used', field: 'Used', sortable: true },
]

function loadAlerts() {
  Keller_bot.getAlerts().then(result => {
    if (result.status === 200) {
      alerts.value = result.data.Data.map(s => KellerBotAlert.fromApi(s))
    }
  }).catch(() => {
    //TODO: Dialog
  })
}

onMounted(() => {
  loadAlerts();
})
</script>

<template>
  <q-page padding>
    <q-table :rows="alerts" :columns="columns" :pagination="emptyPagination" hide-pagination/>
  </q-page>
</template>

<style scoped>

</style>
