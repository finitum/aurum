<template>
  <div class="w-full bg-gray-300 rounded p-4 ">
      {{user}}
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import {client} from "@/client/client";

export default defineComponent({
  name: "ModifyUser",
  async setup() {
    const user = await client.GetUser();

    return {
      user
    };
  }
});
</script>

<style scoped lang="postcss"></style>
