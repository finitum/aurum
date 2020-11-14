<template>
  <div
    class="w-full bg-gray-200 rounded py-8 px-10 mt-3 shadow shadow-md border-1 border-gray-800"
  >
    <form @submit.prevent="" class="grid grid-cols-5 gap-2">
      <h1 class="col-span-full text-center font-semibold text-3xl mb-4">
        {{ application.name }}
        <span v-show="isAdmin(application)">(Admin)</span>
      </h1>
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, onMounted } from "vue";
import { ApplicationWithRole, Role } from "aurum-client";

export default defineComponent({
  name: "Application",
  props: {
    application: {
      type: Object as PropType<ApplicationWithRole>,
      required: true
    }
  },
  setup() {
    function isAdmin(app: ApplicationWithRole): boolean {
      return app.role === Role.Admin;
    }

    return {
      isAdmin
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
