import { LoginForm } from "@/components/forms/login-form";
import { PageLayout } from "@/components/layouts/page-layout";

export default function LoginPage() {
  return (
    <PageLayout
      heading="Welcome Back 🧃"
      subheading="Sign in to continue tracking your drinks"
      className="sm:w-[350px]"
    >
      <LoginForm />
    </PageLayout>
  );
}
