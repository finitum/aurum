<template>
  <div class="flex flex-col w-100 items-center align-middle">
    <div class="lg:w-3/4 w-full p-2">
      <Suspense>
        <template #default>
          <ModifyUser />
        </template>
        <template #fallback>
          <div>Loading...</div>
        </template>
      </Suspense>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import ModifyUser from "../components/ModifiyUser.vue";

export default defineComponent({
  name: "Home",
  components: { ModifyUser },
  setup() {
    return {};
  }
});
</script>
