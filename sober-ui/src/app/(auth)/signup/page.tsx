import { SignupForm } from "@/components/forms/signup-form";

export default function SignupPage() {
  return (
    <div className="container flex h-screen w-full flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold tracking-tight">Sōber 🧃</h1>
          <h2 className="text-lg font-semibold tracking-tight">
            Create your account
          </h2>
          <p className="text-sm text-muted-foreground">
            Enter your details to get started
          </p>
        </div>
        <SignupForm />
      </div>
    </div>
  );
}
