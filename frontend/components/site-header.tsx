"use client";

import Link from "next/link";
import { Briefcase, LogOut } from "lucide-react";

import { useAuth } from "@/components/auth-provider";
import { useSocket } from "@/components/socket-provider";
import { Button } from "@/components/ui/button";

export function SiteHeader() {
  const { user, logout } = useAuth();
  const { connected } = useSocket();

  return (
    <header className="bg-background/80 sticky top-0 z-10 border-b backdrop-blur">
      <div className="mx-auto flex w-full max-w-5xl items-center gap-3 px-4 py-3 sm:px-6">
        <Link href="/" className="flex items-center gap-2 font-semibold">
          <Briefcase className="size-5" aria-hidden />
          Referral Hub
        </Link>

        <div className="ml-auto flex items-center gap-3">
          {user ? (
            <>
              <span
                className="text-muted-foreground flex items-center gap-1.5 text-xs"
                title={connected ? "Live updates connected" : "Reconnecting"}
              >
                <span
                  className={`size-2 rounded-full ${
                    connected ? "bg-emerald-500" : "bg-amber-500"
                  }`}
                  aria-hidden
                />
                <span className="hidden sm:inline">
                  {connected ? "Live" : "Reconnecting"}
                </span>
              </span>

              <span className="hidden text-sm sm:inline">{user.name}</span>

              <Button variant="ghost" size="sm" onClick={logout}>
                <LogOut className="size-4" aria-hidden />
                <span className="sr-only sm:not-sr-only">Sign out</span>
              </Button>
            </>
          ) : (
            <span className="text-muted-foreground text-sm">Not signed in</span>
          )}
        </div>
      </div>
    </header>
  );
}
