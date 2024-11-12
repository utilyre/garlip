"use client";

import { FormEvent, useState } from "react";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      const response = await fetch("/api/v1/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      if (response.status == 422) {
        const error = await response.json();
        switch (error.field) {
          case "username":
            console.log("username error:", error.message);
            break;
          case "password":
            console.log("password error:", error.message);
            break;
        }
      }

      if (!response.ok) {
        throw new Error(`http failure with status ${response.status}`);
      }
    } catch (error) {
      console.error("fetch failed due to", error);
    }
  }

  return (
    <main className="flex h-screen flex-col items-center justify-center p-4">
      <div className="flex w-full flex-col gap-10 sm:w-1/2 md:w-1/3">
        <h1 className="text-center text-2xl font-bold sm:text-4xl">
          Login to Garlip
        </h1>

        <form onSubmit={onSubmit} className="flex flex-col items-center gap-4">
          <div className="flex w-full flex-col gap-2">
            <div className="flex flex-col gap-1">
              <label htmlFor="username">Username</label>
              <input
                id="username"
                type="text"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="h-12 rounded-md border border-gray-500 bg-background p-2 text-sm outline-none transition hover:border-gray-800 focus:border-gray-400 focus:ring focus:ring-gray-300 dark:hover:border-gray-400 dark:focus:ring-gray-600"
              />
            </div>

            <div className="flex flex-col gap-1">
              <label htmlFor="password">Password</label>
              <input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="h-12 rounded-md border border-gray-500 bg-background p-2 text-sm outline-none transition hover:border-gray-800 focus:border-gray-400 focus:ring focus:ring-gray-300 dark:hover:border-gray-400 dark:focus:ring-gray-600"
              />
            </div>
          </div>

          <button
            type="submit"
            className="h-12 w-full rounded-md bg-foreground text-background outline-none transition-colors hover:bg-gray-800 focus:outline-blue-400 dark:hover:bg-gray-300"
          >
            Login
          </button>
        </form>
      </div>
    </main>
  );
}
