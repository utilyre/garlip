"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
      <div className="flex w-full flex-col gap-10 sm:w-1/2 md:w-1/4">
        <h1 className="text-center text-2xl font-bold sm:text-4xl">
          Login to Garlip
        </h1>

        <form onSubmit={onSubmit} className="space-y-4">
          <div className="flex w-full flex-col gap-2">
            <div className="space-y-1">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                type="text"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </div>

            <div className="space-y-1">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
          </div>

          <Button type="submit" className="w-full">
            Login
          </Button>
        </form>
      </div>
    </main>
  );
}
