<script setup lang="ts">
import type {
  ApiKellerBotActionTypeGif
} from 'src/models/keller_bot_action';
import {
  KellerBotAction
} from 'src/models/keller_bot_action';
import { computed, onMounted } from 'vue';
import { api } from 'boot/axios';

const action = defineModel({ type: KellerBotAction, required: true });
const actionGif = computed(() => {
  return action.value.actionGif;
});

const durationSeconds = computed({
  get() {
    return actionGif.value.DurationMs / 1000;
  },
  set(val: number) {
    actionGif.value.DurationMs = val * 1000;
  },
});

onMounted(() => {
  if (action.value.isNew) {
    action.value.Data = {
      DurationMs: 0,
    } as ApiKellerBotActionTypeGif;
  }
});
</script>

<template>
  <div class="row q-mt-lg" v-if="actionGif">
    <div class="col q-pr-lg">
      <q-card>
        <q-card-section class="text-white text-subtitle2 bg-primary">
          Settings
        </q-card-section>
        <q-card-section>
        <q-list dense>
          <q-item>
            <q-item-section avatar>
              <q-icon color="teal" name="sym_o_timer" />
            </q-item-section>
            <q-item-section>
              <q-slider
                class="q-mt-lg"
                :label-value="`Duration: ${durationSeconds}s`"
                label
                label-always
                v-model="durationSeconds"
                :min="0.0"
                :step="0.1"
              />
            </q-item-section>
          </q-item>
        </q-list>
        </q-card-section>
      </q-card>
    </div>
    <div class="col-4">
      <q-card>
        <q-card-section class="text-white text-subtitle2 bg-primary">
          Preview
        </q-card-section>
        <q-card-section class="q-pa-none">
      <q-img
        v-if="actionGif.FileName"
        :src="api.getUri({ method: 'get', url: `/action/${action.ID}` })"
      />
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<style scoped></style>
