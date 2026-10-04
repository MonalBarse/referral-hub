"use client";

import { useState } from "react";
import { toast } from "sonner";

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
import { Textarea } from "@/components/ui/textarea";
import { ApiError, api } from "@/lib/api";

export function ReferralForm({ jobId }: { jobId: string }) {
  const { token } = useAuth();
  const [candidateName, setCandidateName] = useState("");
  const [candidateEmail, setCandidateEmail] = useState("");
  const [message, setMessage] = useState("");
  const [fields, setFields] = useState<Record<string, string>>({});
  const [busy, setBusy] = useState(false);

  async function onSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (!token) return;

    setBusy(true);
    setFields({});
    try {
      await api.createReferral(token, jobId, {
        candidateName,
        candidateEmail,
        message,
      });
      setCandidateName("");
      setCandidateEmail("");
      setMessage("");
      toast.success("Referral submitted");
    } catch (err) {
      if (err instanceof ApiError) {
        setFields(err.fields ?? {});
        // A 409 means this candidate is already referred for this job.
        if (!err.fields) toast.error(err.message);
      } else {
        toast.error("Could not submit the referral");
      }
    } finally {
      setBusy(false);
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Refer a candidate</CardTitle>
        <CardDescription>
          Each candidate can only be referred once per job.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={onSubmit} className="grid gap-4" noValidate>
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="candidate-name">Candidate name</Label>
              <Input
                id="candidate-name"
                value={candidateName}
                onChange={(e) => setCandidateName(e.target.value)}
                placeholder="Priya Nair"
                aria-invalid={Boolean(fields.candidateName)}
              />
              {fields.candidateName && (
                <p className="text-destructive text-sm">
                  Candidate name {fields.candidateName}
                </p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="candidate-email">Candidate email</Label>
              <Input
                id="candidate-email"
                type="email"
                value={candidateEmail}
                onChange={(e) => setCandidateEmail(e.target.value)}
                placeholder="priya@example.com"
                aria-invalid={Boolean(fields.candidateEmail)}
              />
              {fields.candidateEmail && (
                <p className="text-destructive text-sm">
                  Candidate email {fields.candidateEmail}
                </p>
              )}
            </div>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="candidate-message">Why they are a good fit</Label>
            <Textarea
              id="candidate-message"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Optional context for whoever reviews this referral."
              rows={3}
              aria-invalid={Boolean(fields.message)}
            />
            {fields.message && (
              <p className="text-destructive text-sm">Message {fields.message}</p>
            )}
          </div>

          <Button type="submit" disabled={busy} className="justify-self-start">
            {busy ? "Submitting" : "Submit referral"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
