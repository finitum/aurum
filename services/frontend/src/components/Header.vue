<template>
  <header class="flex justify-between flex-wrap bg-primary items-stretch">
    <h1 class="font-semibold text-white m-3">Aurum</h1>

    <div class="flex flex-row justify-center">
      <div
        class="bg-primary cursor-pointer flex flex-col items-center justify-center h-full"
        @click="logoutButton()"
      >
        <span class="button">Log out</span>
      </div>
    </div>
  </header>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { client } from "../client/client";

export default defineComponent({
  name: "Header",
  setup() {
    function logoutButton() {
      client.Logout();
    }

    return {
      logoutButton
    };
  }
});
</script>

<style lang="postcss" scoped>
.button {
  @apply text-sm text-white px-5;
}
</style>
