<template>
  <div>
    <Application
      v-for="(_, index) in applications"
      :key="index"
      v-model:application="applications[index]"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import { ApplicationWithRole } from "aurum-client";
import { client } from "../client/client";
import { CreateNotification } from "./modals/NotificationState";
import Application from "./Application.vue";

export default defineComponent({
  name: "ApplicationList",
  components: { Application },
  async setup() {
    const applications = reactive<ApplicationWithRole[]>([]);

    const apps = await client.GetApplicationsForUser();
    if (apps.isOk()) {
      applications.splice(0, applications.length);
      applications.push(...apps.value);
    } else {
      CreateNotification(apps.error.Message);
    }

    return {
      applications
    };
  }
});
</script>

<style lang="postcss" scoped>
input {
  @apply bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight rounded;
}

.inputEnabled {
  @apply outline-none bg-white border-indigo-500;
}

button {
  @apply px-2 py-2 bg-primary text-white mx-2 rounded;
}
</style>
