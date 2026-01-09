"use client";

import { useEffect } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";

export default function Page() {
  const { loading, authenticated } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading) {
      if (authenticated) {
        router.push("/chat");
      } else {
        router.push("/login");
      }
    }
  }, [loading, authenticated, router]);

  return <div className="flex items-center justify-center h-screen">Loading...</div>;
}