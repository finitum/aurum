<template>
  <div
    class="w-full bg-gray-200 rounded py-8 px-10 shadow shadow-md border-1 border-gray-800"
  >
    <form @submit.prevent="" class="grid grid-cols-5 gap-2">
      <h1 class="col-span-full text-center font-semibold text-3xl mb-4">
        User Information
      </h1>

      <label class="col-start-1 mr-3 font-semibold ml-auto" for="username"
        >Username:</label
      >
      <input
        class="col-start-2 col-span-3"
        id="username"
        v-model="user.username"
        disabled
      />

      <label class="col-start-1 mr-3 font-semibold ml-auto" for="email"
        >Email:</label
      >
      <input
        class="col-start-2 col-span-3"
        id="email"
        v-model="user.email"
        :disabled="!changeEmail"
        :class="{ inputEnabled: changeEmail }"
        tabindex=4
      />
      <button class="col-start-5" @click="doChangeEmail()" tabindex=1>
        <span v-if="!changeEmail">Change</span>
        <span v-if="changeEmail">Update</span>
      </button>

      <label class="mr-3 font-semibold col-start-1 ml-auto" for="password"
        >Password:</label
      >
      <input
        class="col-start-2 col-span-3"
        id="password"
        v-model="user.password"
        placeholder="********"
        :disabled="!changePassword"
        :class="{ inputEnabled: changePassword }"
        type="password"

        tabindex=1
      />
      <button class="col-start-5" @click="doChangePassword()" tabindex=3>
        <span v-if="!changePassword">Change</span>
        <span v-if="changePassword">Update</span>
      </button>

      <input
        class="col-start-2 col-span-3"
        id="repeat"
        v-model="passwordRepeat"
        placeholder="Repeat password"
        v-if="changePassword"
        :class="{ inputEnabled: changePassword }"
        type="password"
        tabindex=2
      />
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, onUnmounted } from "vue";
import { client } from "../client/client";
import { User } from "aurum-client";
import { CreateNotification } from "./modals/NotificationState";

export default defineComponent({
  name: "ModifyUser",
  async setup() {
    const user = ref<User>({} as User);
    const passwordRepeat = ref("");
    const changeEmail = ref(false);
    const changePassword = ref(false);

    onMounted(async () => {
      const userOrError = await client.GetUserInfo();
      if (userOrError.isErr()) {
        const error = userOrError.error;
        CreateNotification(error.Message);
        return;
      } else {
        user.value = userOrError.value;
        user.value.password = "";
      }
    });

    function escapeHandler(e: KeyboardEvent) {
      if (e.key === "Escape") {
        changeEmail.value = false;
        changePassword.value = false;
      }
    }

    onMounted(() => {
      window.addEventListener("keydown", escapeHandler);
    });

    onUnmounted(() => {
      window.removeEventListener("keydown", escapeHandler);
    });

    async function doChangeEmail() {
      if (changeEmail.value) {
        const updateUser: User = {
          email: user.value.email
        } as User;

        const resp = await client.UpdateUser(updateUser);
        if (resp.isOk()) {
          user.value = resp.value;
        } else {
          return CreateNotification(resp.error.Message);
        }
      }

      changeEmail.value = !changeEmail.value;
    }

    async function doChangePassword() {
      if (changePassword.value) {
        if (user.value.password === passwordRepeat.value) {
          const updateUser: User = {
            password: user.value.password
          } as User;

          const resp = await client.UpdateUser(updateUser);
          if (resp.isOk()) {
            user.value = resp.value;
          } else {
            return CreateNotification(resp.error.Message);
          }
        } else {
          return CreateNotification("Passwords don't match");
        }
      }

      changePassword.value = !changePassword.value;
    }

    return {
      user,
      changeEmail,
      changePassword,
      passwordRepeat,
      doChangePassword,
      doChangeEmail
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
