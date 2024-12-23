import { LoginForm } from "@/components/forms/login-form";
import { PageLayout } from "@/components/layouts/page-layout";

export default function LoginPage() {
  return (
    <PageLayout
      heading="Sōber 🧃"
      subheading="Enter your credentials to continue"
      className="sm:w-[350px]"
    >
      <LoginForm />
    </PageLayout>
  );
}
