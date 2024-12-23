import { SignupForm } from "@/components/forms/signup-form";
import { PageLayout } from "@/components/layouts/page-layout";

export default function SignupPage() {
  return (
    <PageLayout
      heading="Create an account"
      subheading="Enter your email below to create your account"
      className="sm:w-[350px]"
    >
      <SignupForm />
    </PageLayout>
  );
}
