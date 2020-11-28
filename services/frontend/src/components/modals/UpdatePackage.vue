<template>
  <div
    class="fixed z-10 inset-0 overflow-y-auto flex justify-center items-center"
  >
    <div class="fixed inset-0 transition-opacity">
      <div
        class="absolute inset-0 bg-gray-500 opacity-75"
        @click="$emit('close')"
      ></div>
    </div>
    <form
      class="bg-white z-20 text-center p-2 my-5 shadow-md rounded grid grid-cols-4 auto-cols-fr gap-8 lg:w-1/2 w-full"
      @submit.prevent="addOrUpdatePackage"
    >
      <h1 class="col-span-full text-3xl font-semibold">{{ pkg.Name }}</h1>

      <label for="name" class="label">Package Name:</label>
      <input id="name" v-model="pkg.Name" disabled class="input" required />

      <label for="repourl" class="label">Repository URL:</label>
      <input
        id="repourl"
        v-model="pkg.RepoURL"
        class="input"
        :disabled="!external"
        required
      />

      <label for="branch" class="label">Repository Branch:</label>
      <input
        id="branch"
        v-model="pkg.RepoBranch"
        :disabled="!external"
        class="input"
        required
      />

      <label for="duration" class="label">Update frequency: </label>
      <DurationPicker
        id="duration"
        class="col-span-3"
        v-model="pkg.UpdateFrequency"
      />

      <label for="lastn" class="label">Keep last:</label>
      <input
        id="lastn"
        type="number"
        v-model="pkg.KeepLastN"
        @focusout="focusOut"
        class="input"
        min="1"
        @focus="$event.target.select()"
        required
      />

      <button class="col-span-full" type="submit" v-if="mode === 'add'">
        Add Package
      </button>
      <button class="col-span-full" type="submit" v-else>Update Package</button>
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref, onMounted, onUnmounted } from "vue";
import DurationPicker from "@/components/DurationPicker.vue";
import {
  NewPackage,
  Package
} from "../../../../../clients/typescript/src/Models";
import { AddPackage, UpdatePackage } from "@/client/API";
import { packages } from "@/client/packages";

export default defineComponent({
  name: "UpdatePackage",
  components: { DurationPicker },

  props: {
    pkgprop: {
      type: Object as PropType<Package>,
      required: true
    },
    mode: {
      type: String as PropType<"update" | "add">,
      required: true
    },
    external: {
      type: Boolean,
      default: false
    }
  },

  setup(props, { emit }) {
    const pkg = ref(NewPackage());

    function escapeHandler(e: KeyboardEvent) {
      if (e.key === "Escape") {
        emit("close");
      }
    }

    onMounted(() => {
      pkg.value = props.pkgprop;
      window.addEventListener("keydown", escapeHandler);
    });

    onUnmounted(() => {
      window.removeEventListener("keydown", escapeHandler);
    });

    function focusOut(e: Event) {
      if ((e.target as HTMLInputElement).value === "") {
        (e.target as HTMLInputElement).value = "1";
      }
    }

    function addOrUpdatePackage() {
      pkg.value.KeepLastN = Number(pkg.value.KeepLastN);

      if (props.mode === "update") {
        UpdatePackage(pkg.value).then(() => emit("close"));
      } else {
        AddPackage(pkg.value).then(() => {
          emit("close");
          packages.push(pkg.value);
        });
      }
    }

    return {
      pkg,
      focusOut,
      addOrUpdatePackage
    };
  }
});
</script>

<style lang="postcss" scoped>
button {
  @apply flex-shrink-0 bg-secondary text-sm text-white py-2 px-3 rounded;
}

.input {
  @apply bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight col-span-3;

  &:focus {
    @apply outline-none bg-white border-secondary;
  }
}

.label {
  @apply col-span-1 text-center block pt-1;
}
</style>
