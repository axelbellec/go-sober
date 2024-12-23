import { SignupForm } from "@/components/forms/signup-form";
import { PageLayout } from "@/components/layouts/page-layout";

export default function SignupPage() {
  return (
    <PageLayout
      heading="Join Sōber 🧃"
      subheading="Create your account to start tracking mindfully"
      className="sm:w-[350px]"
    >
      <SignupForm />
    </PageLayout>
  );
}
