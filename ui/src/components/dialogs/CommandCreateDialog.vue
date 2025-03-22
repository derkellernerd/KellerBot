<script setup lang="ts">
import type {
  ApiKellerBotCommandCreateRequest,
  ApiKellerBotCommandMessageAction,
} from '@/models/keller_bot_command.ts'
import { onMounted, ref } from 'vue'
import Keller_bot from '@/service/keller_bot.ts'
import type { AxiosError } from 'axios'

const commandCreateRequest = ref<ApiKellerBotCommandCreateRequest>()
const messageAction = ref<ApiKellerBotCommandMessageAction>()

const emits = defineEmits(['created'])

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
      emits('created')
      reset()
    }
  }).catch((err : AxiosError) => {
    console.log(err);
  })
}

onMounted(() => {
  reset()
})
</script>

<template>
  <div v-if="commandCreateRequest && messageAction">
    <form class="max-w-sm mx-auto" @submit.prevent="createCommand" @reset="reset">
      <div class="mb-5">
        <label for="command" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Command</label>
        <input id="command" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" v-model="commandCreateRequest.Command" />
      </div>

      <div class="mb-5">
        <input id="type_message_action" type="radio" value="MESSAGE_ACTION" class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600" v-model="commandCreateRequest.Type">
        <label for="type_message_action" class="ms-2 text-sm font-medium text-gray-400 dark:text-gray-500">Message</label>
      </div>

      <div class="mb-5">
        <label for="message" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Message</label>
        <input id="message" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" v-model="messageAction.Message" />
      </div>

      <button type="submit" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800" value="submit">Submit</button>
      <button type="button" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800" value="reset">Reset</button>
    </form>
  </div>
</template>

<style scoped>
</style>
