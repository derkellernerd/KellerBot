<script setup lang="ts">
import ChatEventComponent from 'components/event/ChatEventComponent.vue';
</script>

<template>
  <q-layout>
    <q-page-container>
      <q-page class="bg-black text-white text-h4">
        <ChatEventComponent>
          <template v-slot:body="props">
            <div v-for="(message, index) in props.messages" :key="index">
              <span class="text-bold">{{message.User}}</span> {{message.Message}}
            </div>
          </template>
        </ChatEventComponent>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<style scoped>
</style>
