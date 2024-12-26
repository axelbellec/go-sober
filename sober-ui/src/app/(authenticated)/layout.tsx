"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function AuthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem(
      process.env.NEXT_PUBLIC_LOCALSTORAGE_TOKEN_KEY!
    );
    if (!token) {
      router.push("/login");
    }
  }, [router]);

  return <>{children}</>;
}
