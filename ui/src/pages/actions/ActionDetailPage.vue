<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  type ApiKellerBotAction,
  KellerBotAction,
  KellerBotActionType,
} from 'src/models/keller_bot_action';
import Keller_bot from 'src/service/keller_bot';
import ActionGifEdit from 'components/action/ActionGifEdit.vue';
import ActionCompositionEdit from 'components/action/ActionCompositionEdit.vue';
import ActionSoundEdit from 'components/action/ActionSoundEdit.vue';
import ActionTextEdit from 'components/action/ActionTextEdit.vue';
import ActionChatMessageEdit from 'components/action/ActionChatMessageEdit.vue';
import ActionHttpEdit from 'components/action/ActionHttpEdit.vue';
import { showSuccessToast } from 'src/helper/functions';

const route = useRoute();
const router = useRouter();
const action = ref<KellerBotAction>();
const file = ref<File>();
const fileHasChanged = ref(false);

const actionId = computed(() => {
  return parseInt(route.params.actionId as string);
});

const isNew = computed(() => {
  return route.path.endsWith('create');
});

const actionTypes = computed(() => {
  return KellerBotAction.actionTypes();
});

const isFileAction = computed(() => {
  const possibleTypes = [KellerBotActionType.Gif, KellerBotActionType.Sound];
  return possibleTypes.filter((s) => s == action.value?.ActionType).length > 0;
});

async function uploadFile() {
  if (!file.value) return;
  if (!isFileAction.value) return;
  const actionData = action.value as KellerBotAction;

  await Keller_bot.uploadActionFile(actionData.ID, file.value);
}

function createOrSave() {
  if (isNew.value) {
    create();
  } else {
    save();
  }
}

function save() {
  if (!action.value) return;

  void Keller_bot.updateAction(action.value.ID, action.value)
    .then(async (result) => {
      if (result.status === 201) {
        action.value = KellerBotAction.fromApi(result.data.Data);
        await uploadFile();
        loadAction();
        showSuccessToast(`Action ${action.value.ActionName} saved`)
      }
    })
}

function create() {
  if (!action.value) return;

  void Keller_bot.createAction(action.value)
    .then(async (result) => {
      if (result.status === 201) {
        action.value = KellerBotAction.fromApi(result.data.Data);
        await uploadFile();
        showSuccessToast(`Action ${action.value.ActionName} created`)
        await router.push({ name: 'ActionOverview' });
      }
    })
}

function loadAction() {
  Keller_bot.getAction(actionId.value)
    .then((result) => {
      if (result.status === 200) {
        action.value = KellerBotAction.fromApi(result.data.Data);
      }
    })
    .catch(() => {
      //TODO: catch
    });
}

onMounted(() => {
  if (isNew.value) {
    action.value = KellerBotAction.fromApi({} as ApiKellerBotAction);
  } else {
    loadAction();
  }
});
</script>

<template>
  <q-page padding>
    <div v-if="action">
      <q-form @submit="createOrSave()">
        <q-card>
          <q-card-section class="text-h6 bg-primary text-white"> Common</q-card-section>
          <q-card-section>
            <q-input v-model="action.ActionName" label="Name" :readonly="!isNew" />
            <q-select
              v-model="action.ActionType"
              label="Action Type"
              :readonly="!isNew"
              :options="actionTypes"
            />
            <q-file
              label="File"
              v-if="isFileAction"
              v-model="file"
              @change="fileHasChanged = true"
            />
          </q-card-section>
        </q-card>
        <ActionGifEdit v-if="action.ActionType == KellerBotActionType.Gif" v-model="action" />
        <ActionSoundEdit v-if="action.ActionType == KellerBotActionType.Sound" v-model="action" />
        <ActionTextEdit v-if="action.ActionType == KellerBotActionType.Text" v-model="action" />
        <ActionChatMessageEdit
          v-if="action.ActionType == KellerBotActionType.ChatMessage"
          v-model="action"
        />
        <ActionCompositionEdit
          v-if="action.ActionType == KellerBotActionType.Composition"
          v-model="action"
        />
        <ActionHttpEdit v-if="action.ActionType == KellerBotActionType.Http" v-model="action" />
        <q-btn
          type="submit"
          :icon="isNew ? 'add' : 'save'"
          :label="isNew ? 'create' : 'save'"
          class="full-width q-mt-lg"
          color="positive"
        />
      </q-form>
    </div>
  </q-page>
</template>

<style scoped></style>
