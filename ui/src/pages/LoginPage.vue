<script setup lang="ts">
import Keller_bot from 'src/service/keller_bot';

function login() {
  Keller_bot.twitchLogin()
    .then((result) => {
      console.log(result);
    })
    .catch((err) => {
      console.log(err);
    });
}
</script>

<template>
  <q-layout>
    <q-page-container>
      <q-page padding>
        <q-form @submit="login">
          <q-btn icon="add" type="submit" />
        </q-form>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<style scoped></style>
