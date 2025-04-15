<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title @click="gotoDashboard">
          KellerBot
        </q-toolbar-title>

        <div>v0.1.0</div>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
    >
      <q-list>
        <q-item-label
          header
        >
          Menu
        </q-item-label>
        <q-item clickable v-ripple :to="{name: 'ChatCommandOverview'}">
          <q-item-section>
            <q-item-label>Command</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable v-ripple :to="{name: 'ActionOverview'}">
          <q-item-section>
            <q-item-label>Actions</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable v-ripple :to="{name: 'Event'}">
          <q-item-section>
            <q-item-label>Event</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable v-ripple :to="{name: 'TwitchEvent'}">
          <q-item-section>
            <q-item-label>Twitch Events</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const leftDrawerOpen = ref(false);
const router = useRouter();

function toggleLeftDrawer () {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}

function gotoDashboard() {
  void router.push({name: 'Dashboard'})
}
</script>
