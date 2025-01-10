import { ProfileForm } from "@/components/forms/profile-form";
import { PageLayout } from "@/components/layouts/page-layout";

export default function ProfilePage() {
  return (
    <PageLayout
      heading="Profile Settings"
      subheading="Update your personal information"
      className="sm:w-[350px]"
    >
      <ProfileForm />
    </PageLayout>
  );
}
