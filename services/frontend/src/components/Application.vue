<template>
  <div class="w-full bg-gray-300 rounded py-8 px-10 ">
    <form @submit.prevent="" class="grid grid-cols-5 gap-2">
      <h1 class="col-span-full text-center font-semibold text-3xl mb-4">User Information</h1>




    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, onUnmounted } from "vue";
import { client } from "@/client/client";
import { AurumError, ErrorCode, User } from "aurum-client";
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
      if (userOrError.IsError()) {
        const error = userOrError.GetError() as AurumError;
        CreateNotification(error.Message);
        return;
      } else {
        user.value = userOrError.GetOk() as User;
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
        if (resp.IsOk()) {
          user.value = resp.GetOk() as User;
        } else {
          return CreateNotification((resp.GetError() as AurumError).Message);
        }
      }

      changeEmail.value = !changeEmail.value;
    }

    async function doChangePassword() {
      if (changePassword.value) {
        if (user.value.password === passwordRepeat.value){
          const updateUser: User = {
            password: user.value.password
          } as User;

          const resp = await client.UpdateUser(updateUser);
          if (resp.IsOk()) {
            user.value = resp.GetOk() as User;
          } else {
            return CreateNotification((resp.GetError() as AurumError).Message);
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
