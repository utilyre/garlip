export default function Login() {
  return (
    <main className="flex h-screen flex-col items-center justify-center p-4">
      <div className="w-full border-2 border-accent-700 bg-accent-50 p-4 sm:w-1/2 md:w-1/3">
        <h1 className="text-center text-2xl font-light text-foreground">
          Who are you?
        </h1>

        <div className="h-6"></div>

        <form className="flex flex-col items-center gap-6">
          <div className="flex w-full flex-col gap-4">
            <div className="flex flex-col">
              <label htmlFor="username" className="text-xs">
                Username
              </label>
              <input
                id="username"
                type="text"
                placeholder="donaldtrump"
                className="h-8 border border-foreground p-2 text-sm outline-none"
              />
            </div>

            <div className="flex flex-col">
              <label htmlFor="password" className="text-xs">
                Password
              </label>
              <input
                id="password"
                type="password"
                placeholder="n#tgAL7m"
                className="h-8 border border-foreground p-2 text-sm outline-none"
              />
            </div>
          </div>

          <button
            type="submit"
            className="h-10 w-40 bg-accent-800 text-sm font-bold text-background outline-none focus:ring-2 focus:ring-accent-800 focus:ring-offset-2"
          >
            Login
          </button>
        </form>
      </div>
    </main>
  );
}
