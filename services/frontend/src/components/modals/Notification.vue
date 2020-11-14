<template>
  <div
    v-if="notificationState.enabled"
    class="fixed-z z-10 absolute right-0 bottom-0 shadow-md p-3 m-5 rounded"
    v-bind:style="{ backgroundColor: notificationState.color }"
  >
    <div class="absolute cursor-pointer right-0 top-0 m-2" @click="close()">
      <svg
        class="fill-current text-black"
        xmlns="http://www.w3.org/2000/svg"
        width="18"
        height="18"
        viewBox="0 0 18 18"
      >
        <path
          d="M14.53 4.53l-1.06-1.06L9 7.94 4.53 3.47 3.47 4.53 7.94 9l-4.47 4.47 1.06 1.06L9 10.06l4.47 4.47 1.06-1.06L10.06 9z"
        ></path>
      </svg>
    </div>
    <div class="p-2 mx-2">
      {{ notificationState.message }}
    </div>
  </div>
</template>

<script lang="typescript">
import { defineComponent } from "vue";
import { notificationState } from "@/components/modals/NotificationState";

export default defineComponent({
  name: "Notification",
  setup() {
    function close() {
      notificationState.enabled = false;
    }

    return {
      notificationState,
      close
    };
  }
});
</script>
