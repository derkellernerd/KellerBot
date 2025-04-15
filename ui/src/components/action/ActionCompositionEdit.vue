<script setup lang="ts">
import { KellerBotAction } from 'src/models/keller_bot_action';
import type { ApiKellerBotActionTypeComposition } from 'src/models/keller_bot_action';
import { computed, onMounted, ref } from 'vue';
import Keller_bot from 'src/service/keller_bot';

const action = defineModel({ type: KellerBotAction, required: true });
const actionComposition = computed(() => {
  return action.value.actionComposition;
});

const actions = ref<KellerBotAction[]>([]);

function loadActions() {
  Keller_bot.getActions()
    .then((result) => {
      if (result.status === 200) {
        actions.value = result.data.Data.map((s) => KellerBotAction.fromApi(s));
      }
    })
    .catch(() => {
      //TODO: catch
    });
}

const durationSeconds = computed({
  get() {
    return actionComposition.value.DurationMs / 1000;
  },
  set(val: number) {
    actionComposition.value.DurationMs = val * 1000;
  },
});

onMounted(() => {
  console.log('mounted: ', action.value);
  if (action.value.isNew) {
    console.log('new');
    action.value.Data = {
      Actions: [''] as string[],
      PostAction: '',
    } as ApiKellerBotActionTypeComposition;
  }

  loadActions();
});
</script>

<template>
  <q-card class="q-mt-lg" v-if="actionComposition">
    <q-card-section class="bg-primary text-white text-h6">
      Actions
    </q-card-section>
    <q-card-section>
      <q-list>
        <q-item v-for="index in actionComposition.Actions.length" :key="index">
          <q-item-section>
            <q-select
              v-model="actionComposition.Actions[index - 1]"
              :options="actions"
              map-options
              emit-value
              option-label="ActionName"
              option-value="ActionName"
            />
          </q-item-section>
          <q-item-section side>
            <q-btn
              color="negative"
              icon="delete"
              @click="actionComposition.Actions.splice(index, 1)"
            />
          </q-item-section>
        </q-item>
        <q-item>
          <q-item-section>
            <q-btn
              icon="add"
              label="Add"
              color="positive"
              class="full-width"
              @click="actionComposition.Actions.push('')"
            />
          </q-item-section>
        </q-item>
      </q-list>
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
        <q-item>
          <q-item-section>
            <q-select
              label="Post Action"
              v-model="actionComposition.PostAction"
              :options="actions"
              map-options
              emit-value
              option-label="ActionName"
              option-value="ActionName"
            />
          </q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>
</template>

<style scoped></style>
