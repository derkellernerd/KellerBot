<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { KellerBotCommand } from '@/models/keller_bot_command.ts'
import KellerBot from '@/service/keller_bot.ts'
import CommandCreateDialog from '@/components/dialogs/CommandCreateDialog.vue'

const commands = ref<KellerBotCommand[]>([]);

function loadCommands() {
  KellerBot.getCommands().then(result => {
    if (result.status === 200) {
      commands.value = result.data.Data.map(s => KellerBotCommand.fromApi(s))
    }
  })
}

onMounted(() => {
  loadCommands();
})
</script>

<template>
<div>
  <CommandCreateDialog @created="loadCommands"/>
  <table>
    <tr v-for="command in commands" :key="command.Id">
      <td>{{command.Command}}</td>
      <td>{{command.Type}}</td>
    </tr>
  </table>
</div>
</template>

<style scoped>

</style>
