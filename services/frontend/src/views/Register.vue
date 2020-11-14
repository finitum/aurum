<template>
  <div
    class="fixed z-10 inset-0 overflow-y-auto flex justify-center items-center"
  >
    <div class="fixed inset-0 transition-opacity">
      <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
    </div>
    <div class="bg-white z-20 text-center p-2 my-2 shadow-md rounded lg:w-1/2 w-11/12" style="width: 10cm">
      <h3 class="text-2xl font-medium">Register</h3>
      <form
        class="flex flex-col items-center align-middle justify-between text-center"
        @submit.prevent="ClickRegister"
      >
        <span class="label">
          <label for="username">Username:</label>
          <input
            v-model="user.username"
            id="username"
            type="text"
            class="input-box"
          />
        </span>
        <span class="flex flex-col label">
          <label for="email">Email address:</label>
          <input
            v-model="user.email"
            id="email"
            class="input-box"
            type="email"
          />
        </span>
        <span class="flex flex-col label">
          <label for="password">Password:</label>
          <input
            v-model="user.password"
            id="password"
            type="password"
            class="input-box"
          />
        </span>
        <span class="flex flex-col label">
          <label for="password">Repeat password:</label>
          <input
            v-model="passwordRepeat"
            id="repeatPassword"
            type="password"
            class="input-box"
          />
        </span>

        <button class="button w-full" type="submit">
          Register
        </button>

        <span class="mt-3 text-gray-600">
          Already have an account?
          <a class="text-blue-800 cursor-pointer" @click="router.push('/login')">
            Login
          </a>
          here
        </span>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import router from "../router";
import { client } from "@/client/client";
import { User, ErrorCode } from "aurum-client";
import { CreateNotification } from "@/components/modals/NotificationState";

export default defineComponent({
  name: "Login",
  setup() {
    const user = reactive<User>({ username: "", password: "", email: "" });
    const passwordRepeat = ref("");

    async function ClickRegister() {
      if (user.username === "") {
        CreateNotification("Username field empty");
      } else if (user.email === "") {
        CreateNotification("Email field empty");
      } else if (user.password === "") {
        CreateNotification("Password field empty");
      } else if (user.password !== passwordRepeat.value) {
        CreateNotification("Passwords don't match");
      } else {
        const error = await client.Register(user);
        if (error.isOk()) {
          await router.push("/");
        } else if (error.error.Code === ErrorCode.Unauthorized) {
          CreateNotification("Username or password incorrect");
        } else if (error.error.Code === ErrorCode.WeakPassword) {
          CreateNotification("Password too weak");
        } else {
          CreateNotification(error.error.Message);
        }
      }
    }

    return {
      user,
      ClickRegister,
      passwordRepeat,
      router
    };
  }
});
</script>

<style scoped lang="postcss">
.label {
  @apply flex flex-col text-gray-700 font-bold p-2 m-1 w-full;
}

.button {
  @apply flex-shrink-0 bg-primary text-sm text-white py-2 px-3 rounded;
}

.button:hover {
  @apply bg-primarydark;
}

.input-box {
  @apply bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight;
}

.input-box:focus {
  @apply outline-none bg-white border-indigo-500;
}
</style>
