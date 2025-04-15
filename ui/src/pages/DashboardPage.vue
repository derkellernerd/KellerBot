<script setup lang="ts">

import { KellerBotAction } from 'src/models/keller_bot_action';
import { KellerBotChatCommand } from 'src/models/keller_bot_chat_command';
import { onMounted, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';

const actions = ref<KellerBotAction[]>([]);
const commands = ref<KellerBotChatCommand[]>([]);

function loadCommands() {
  void Keller_bot.getChatCommands().then(result => {
    commands.value = result.data.Data.map(s => KellerBotChatCommand.fromApi(s))
  })
}

function loadActions() {
  void Keller_bot.getActions().then(result => {
    actions.value = result.data.Data.map(s => KellerBotAction.fromApi(s))
  })
}

onMounted(() => {
  loadCommands()
  loadActions()
})
</script>

<template>
<q-page padding>
  <div class="row">
  <q-card class="col-6">
    <q-card-section class="text-h6 text-white bg-primary">
      Stats
    </q-card-section>
    <q-card-section>
      <q-list>
        <q-item>
          <q-item-section>
            Actions
          </q-item-section>
          <q-item-section side>
            {{ actions.length }}
          </q-item-section>
        </q-item>
        <q-item>
          <q-item-section>
            Chat Commands
          </q-item-section>
          <q-item-section side>
            {{ commands.length }}
          </q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>
  </div>
</q-page>
</template>

<style scoped>

</style>
