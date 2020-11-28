<template>
  <div>
    <div
      class="fixed z-10 inset-0 overflow-y-auto flex justify-center items-center"
    >
      <div class="fixed inset-0 transition-opacity">
        <div
          class="absolute inset-0 bg-gray-500 opacity-75"
          @click="$emit('close')"
        ></div>
      </div>
      <div
        class="bg-white z-20 text-center py-3 px-6 my-5 shadow-md rounded grid grid-cols-4 auto-cols-fr gap-8"
      >
        <h1 class="col-span-full text-2xl font-semibold border-b-2">
          <slot name="header"></slot>
        </h1>

        <p class="col-span-full">
          <slot name="message"></slot>
        </p>

        <button v-if="mode === 'Info'" class="info" @click="$emit('close')">
          Ok
        </button>

        <button
          v-if="mode === 'Confirm'"
          class="accept"
          @click="$emit('accept')"
        >
          Accept
        </button>
        <button
          v-if="mode === 'Confirm'"
          class="cancel"
          @click="$emit('close')"
        >
          Cancel
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, onUnmounted, PropType } from "vue";

export default defineComponent({
  name: "Dialog",
  components: {},
  props: {
    mode: {
      type: String as PropType<"Confirm" | "Info">,
      required: true
    }
  },
  setup(_, { emit }) {
    function escapeHandler(e: KeyboardEvent) {
      if (e.key === "Escape") {
        emit("close");
      }
    }

    onMounted(() => {
      window.addEventListener("keydown", escapeHandler);
    });

    onUnmounted(() => {
      window.removeEventListener("keydown", escapeHandler);
    });

    return {};
  }
});
</script>

<style scoped>
button {
  @apply flex-shrink-0 bg-secondary text-sm text-white py-2 px-3 rounded;
}

.info {
  @apply col-span-full bg-info;
}

.accept,
.cancel {
  @apply col-span-2;
}

.accept {
  @apply bg-accept;
}
.cancel {
  @apply bg-cancel;
}
</style>
