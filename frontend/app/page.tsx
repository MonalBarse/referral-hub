"use client";

import { useCallback, useEffect, useState } from "react";

import { useAuth } from "@/components/auth-provider";
import { JobCard } from "@/components/job-card";
import { JobForm } from "@/components/job-form";
import { LoginForm } from "@/components/login-form";
import { useTopic } from "@/components/socket-provider";
import { Skeleton } from "@/components/ui/skeleton";
import { api } from "@/lib/api";
import type { Job, WsEvent } from "@/lib/types";

export default function DashboardPage() {
  const { user, ready } = useAuth();
  const [jobs, setJobs] = useState<Job[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api
      .listJobs()
      .then((res) => setJobs(res.jobs))
      .catch(() => setError("Could not load the job feed. Is the backend running?"));
  }, []);

  const onFeedEvent = useCallback((event: WsEvent) => {
    if (event.type === "job.created") {
      const job = event.payload;
      setJobs((current) => {
        if (!current) return current;
        if (current.some((j) => j.id === job.id)) return current;
        return [job, ...current];
      });
      return;
    }

    if (event.type === "referral.created") {
      const referral = event.payload;
      setJobs((current) =>
        current
          ? current.map((job) =>
              job.id === referral.jobId
                ? { ...job, referralCount: job.referralCount + 1 }
                : job,
            )
          : current,
      );
    }
  }, []);

  useTopic("jobs", onFeedEvent);

  return (
    <div className="grid gap-6">
      {ready && (user ? <JobForm /> : <LoginForm />)}

      <section className="grid gap-3">
        <h2 className="text-lg font-semibold">Active job listings</h2>

        {error && <p className="text-destructive text-sm">{error}</p>}

        {!jobs && !error && (
          <div className="grid gap-3">
            <Skeleton className="h-28 w-full" />
            <Skeleton className="h-28 w-full" />
          </div>
        )}

        {jobs?.length === 0 && (
          <p className="text-muted-foreground text-sm">
            No jobs yet. Post the first one above.
          </p>
        )}

        {jobs?.map((job) => <JobCard key={job.id} job={job} />)}
      </section>
    </div>
  );
}
