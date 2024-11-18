"use client";

import LoginForm from "@/components/login-form";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Login() {
  const router = useRouter();

  useEffect(() => {
    (async () => {
      try {
        const response = await fetch("/api/v1/auth/check");
        if (response.status >= 200 && response.status < 300) {
          router.replace("/");
        }
        if (response.status !== 401 && !response.ok) {
          throw new Error(`http failure with status ${response.status}`);
        }
      } catch (error) {
        console.error("fetch failed due to", error);
      }
    })();
  }, []);

  return (
    <main className="flex h-screen flex-col items-center justify-center p-4">
      <LoginForm className="w-full sm:w-1/2 md:w-1/4" />
    </main>
  );
}
