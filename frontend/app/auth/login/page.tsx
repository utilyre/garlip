import LoginForm from "@/components/login-form";

export default function Login() {
  return (
    <main className="flex h-screen flex-col items-center justify-center p-4">
      <LoginForm className="w-full sm:w-1/2 md:w-1/4" />
    </main>
  );
}
