import { LoginForm } from "@/components/forms/login-form";

export default function LoginPage() {
  return (
    <div className="container flex h-screen w-full flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">Sōber 🧃</h1>
          <h2 className="text-lg font-semibold tracking-tight">
            <span className="text-lg font-semibold">
              Secure Your Alcohol Consumption Tracking
            </span>
          </h2>
          <p className="text-sm text-muted-foreground">
            Enter your credentials to continue
          </p>
        </div>
        <LoginForm />
      </div>
    </div>
  );
}
