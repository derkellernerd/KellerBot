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

        <q-toolbar-title>
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
        <q-item clickable v-ripple :to="{name: 'CommandOverview'}">
          <q-item-section>
            <q-item-label>Command</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable v-ripple :to="{name: 'AlertOverview'}">
          <q-item-section>
            <q-item-label>Alert</q-item-label>
          </q-item-section>
        </q-item>
        <q-item clickable v-ripple :to="{name: 'Event'}">
          <q-item-section>
            <q-item-label>Event</q-item-label>
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

const leftDrawerOpen = ref(false);

function toggleLeftDrawer () {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}
</script>
