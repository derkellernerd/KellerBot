<script setup lang="ts">
import  { KellerBotAction } from 'src/models/keller_bot_action';
import { onMounted, ref } from 'vue';
import { emptyPagination } from 'src/helper/common';
import Keller_bot from 'src/service/keller_bot';

const actions = ref<KellerBotAction[]>([]);

const columns = [
  { name: 'id', label: 'ID', field: 'ID', sortable: true },
  { name: 'actionName', label: 'ActionName', field: 'ActionName', sortable: true },
  { name: 'actionType', label: 'ActionType', field: 'ActionType', sortable: true },
];

function loadActions() {
  Keller_bot.getActions()
    .then((result) => {
      if (result.status === 200) {
        actions.value = result.data.Data.map((s) => KellerBotAction.fromApi(s));
      }
    })
    .catch(() => {
      //TODO: Dialog
    });
}

onMounted(() => {
  loadActions();
});
</script>

<template>
  <q-page padding>
    <q-btn
      class="full-width q-mb-lg"
      color="positive"
      icon="add"
      label="add action"
      :to="{name:'ActionCreate'}"
    />
    <q-table :rows="actions" :columns="columns" :pagination="emptyPagination" hide-pagination>
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th
            v-for="col in props.cols"
            :key="col.name"
            :props="props"
          >
            {{ col.label }}
          </q-th>
          <q-th auto-width />
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td
            v-for="col in props.cols"
            :key="col.name"
            :props="props"
          >
            {{ col.value }}
          </q-td>
          <q-td auto-width>
            <q-btn color="primary" icon="edit" :to="{name: 'ActionDetail', params: {actionId: props.row.ID}}"/>
          </q-td>
        </q-tr>
      </template>

    </q-table>
  </q-page>
</template>

<style scoped></style>
