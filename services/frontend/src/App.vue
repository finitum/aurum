<template>
  <Header />
  <router-view />
  <Notification />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Header from "@/components/Header.vue";
import Notification from "@/components/modals/Notification.vue";

export default defineComponent({
  name: "App",
  components: { Header, Notification },
  setup() {
    return {};
  }
});
</script>

<style lang="postcss">
@tailwind base;
@tailwind components;
@tailwind utilities;

html {
  background: #f4f6f7;
}
</style>
