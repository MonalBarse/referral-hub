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

export function JobForm() {
  const { token } = useAuth();
  const [title, setTitle] = useState("");
  const [company, setCompany] = useState("");
  const [description, setDescription] = useState("");
  const [fields, setFields] = useState<Record<string, string>>({});
  const [busy, setBusy] = useState(false);

  async function onSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (!token) return;

    setBusy(true);
    setFields({});
    try {
      await api.createJob(token, { title, company, description });
      // The new job arrives over the socket like everyone else's, so there is
      // nothing to insert here. Just clear the form.
      setTitle("");
      setCompany("");
      setDescription("");
      toast.success("Job posted");
    } catch (err) {
      if (err instanceof ApiError) {
        setFields(err.fields ?? {});
        if (!err.fields) toast.error(err.message);
      } else {
        toast.error("Could not post the job");
      }
    } finally {
      setBusy(false);
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Post a job</CardTitle>
        <CardDescription>
          It appears in every connected tab immediately.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={onSubmit} className="grid gap-4" noValidate>
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="job-title">Title</Label>
              <Input
                id="job-title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Senior Backend Engineer"
                aria-invalid={Boolean(fields.title)}
              />
              {fields.title && (
                <p className="text-destructive text-sm">Title {fields.title}</p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="job-company">Company</Label>
              <Input
                id="job-company"
                value={company}
                onChange={(e) => setCompany(e.target.value)}
                placeholder="Acme Inc"
                aria-invalid={Boolean(fields.company)}
              />
              {fields.company && (
                <p className="text-destructive text-sm">Company {fields.company}</p>
              )}
            </div>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="job-description">Description</Label>
            <Textarea
              id="job-description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="What the role involves, and what you are looking for."
              rows={4}
              aria-invalid={Boolean(fields.description)}
            />
            {fields.description && (
              <p className="text-destructive text-sm">
                Description {fields.description}
              </p>
            )}
          </div>

          <Button type="submit" disabled={busy} className="justify-self-start">
            {busy ? "Posting" : "Post job"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
