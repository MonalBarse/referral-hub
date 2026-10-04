"use client";

import { useState } from "react";

import { useAuth } from "@/components/auth-provider";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ApiError } from "@/lib/api";

export function LoginForm() {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [fields, setFields] = useState<Record<string, string>>({});
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function onSubmit(event: React.FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError(null);
    setFields({});
    try {
      await login(email, name);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
        setFields(err.fields ?? {});
      } else {
        setError("Something went wrong. Please try again.");
      }
    } finally {
      setBusy(false);
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Sign in</CardTitle>
        <CardDescription>
          Any email works. The first sign in creates the account, so open a second
          browser tab with a different email to see live updates between two users.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={onSubmit} className="grid gap-4" noValidate>
          <div className="grid gap-2">
            <Label htmlFor="login-name">Display name</Label>
            <Input
              id="login-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Alice Chen"
              autoComplete="name"
              aria-invalid={Boolean(fields.name)}
            />
            {fields.name && (
              <p className="text-destructive text-sm">Display name {fields.name}</p>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="login-email">Email</Label>
            <Input
              id="login-email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="alice@example.com"
              autoComplete="email"
              aria-invalid={Boolean(fields.email)}
            />
            {fields.email && (
              <p className="text-destructive text-sm">Email {fields.email}</p>
            )}
          </div>

          {error && <p className="text-destructive text-sm">{error}</p>}

          <Button type="submit" disabled={busy}>
            {busy ? "Signing in" : "Sign in"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
