"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Eye, EyeOff } from "lucide-react";
import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

function capitalizeFirstLetter(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

export default function SignUp() {
  const router = useRouter();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [fullname, setFullname] = useState("");

  const [usernameError, setUsernameError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [fullnameError, setFullnameError] = useState("");
  const [formError, setFormError] = useState("");

  const [showPassword, setShowPassword] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      const signupResponse = await fetch("/api/v1/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password, fullname }),
      });

      if (signupResponse.status === 422) {
        const error = await signupResponse.json();
        switch (error.field) {
          case "username":
            setUsernameError(capitalizeFirstLetter(error.message));
            setPasswordError("");
            setFullnameError("");
            setFormError("");
            break;
          case "password":
            setPasswordError(capitalizeFirstLetter(error.message));
            setUsernameError("");
            setFullnameError("");
            setFormError("");
            break;
          case "fullname":
            setFullnameError(capitalizeFirstLetter(error.message));
            setUsernameError("");
            setPasswordError("");
            setFormError("");
            break;
        }
        return;
      }
      if (signupResponse.status === 409) {
        setFormError("Account already exists");
        setUsernameError("");
        setPasswordError("");
        setFullnameError("");
        return;
      }
      if (!signupResponse.ok) {
        throw new Error(
          `http failure on signup with status ${signupResponse.status}`,
        );
      }

      const loginResponse = await fetch("/api/v1/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      // NOTE: do not handle not found or validation errors since we just signed
      // up with the same information

      if (!loginResponse.ok) {
        throw new Error(
          `http failure on login with status ${signupResponse.status}`,
        );
      }

      router.push("/");
    } catch (error) {
      console.error("fetch failed due to", error);
    }
  }

  return (
    <main className="flex h-screen flex-col items-center justify-center p-4">
      <div className="flex w-full flex-col gap-10 sm:w-1/2 md:w-1/4">
        <h1 className="text-center text-2xl font-bold sm:text-4xl">
          Create an account
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
              {usernameError && (
                <p className="text-sm text-destructive">{usernameError}</p>
              )}
            </div>

            <div className="space-y-1">
              <Label htmlFor="password">Password</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  required
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="absolute right-2 top-1/2 -translate-y-1/2"
                  onClick={() => setShowPassword(!showPassword)}
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </Button>
              </div>
              {passwordError && (
                <p className="text-sm text-destructive">{passwordError}</p>
              )}
            </div>

            <div className="space-y-1">
              <Label htmlFor="fullname">Full name</Label>
              <Input
                id="fullname"
                type="text"
                value={fullname}
                onChange={(e) => setFullname(e.target.value)}
              />
              {fullnameError && (
                <p className="text-sm text-destructive">{fullnameError}</p>
              )}
            </div>
          </div>

          <div className="space-y-1">
            {formError && (
              <p className="text-center text-sm text-destructive">
                {formError}
              </p>
            )}
            <Button type="submit" className="w-full">
              Sign Up
            </Button>
          </div>
        </form>
      </div>
    </main>
  );
}
