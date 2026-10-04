"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { useAuth } from "@/components/auth-provider";
import { ReferralCard } from "@/components/referral-card";
import { ReferralForm } from "@/components/referral-form";
import { useTopic } from "@/components/socket-provider";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { api } from "@/lib/api";
import type { Job, Referral, WsEvent } from "@/lib/types";

export default function JobDetailPage() {
  const params = useParams<{ id: string }>();
  const jobId = params.id;
  const { token } = useAuth();

  const [job, setJob] = useState<Job | null>(null);
  const [referrals, setReferrals] = useState<Referral[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!jobId) return;
    let cancelled = false;

    api
      .getJob(jobId)
      .then((res) => {
        if (cancelled) return;
        setJob(res.job);
        setReferrals(res.referrals);
      })
      .catch(() => {
        if (!cancelled) setError("That job could not be loaded.");
      });

    return () => {
      cancelled = true;
    };
  }, [jobId]);

  const onJobEvent = useCallback((event: WsEvent) => {
    if (event.type !== "referral.created") return;
    const incoming = event.payload;
    setReferrals((current) => {
      if (!current) return current;
      if (current.some((r) => r.id === incoming.id)) return current;
      return [incoming, ...current];
    });
  }, []);

  useTopic(jobId ? `job:${jobId}` : null, onJobEvent);

  if (error) {
    return (
      <div className="grid gap-4">
        <p className="text-destructive text-sm">{error}</p>
        <Link href="/" className="text-sm hover:underline">
          Back to all jobs
        </Link>
      </div>
    );
  }

  return (
    <div className="grid gap-6">
      <Link
        href="/"
        className="text-muted-foreground flex items-center gap-1.5 text-sm hover:underline"
      >
        <ArrowLeft className="size-4" aria-hidden />
        All jobs
      </Link>

      {!job ? (
        <Skeleton className="h-36 w-full" />
      ) : (
        <Card>
          <CardHeader className="gap-1">
            <div className="flex flex-wrap items-start gap-2">
              <CardTitle className="text-xl">{job.title}</CardTitle>
              <Badge className="ml-auto shrink-0">{job.status}</Badge>
            </div>
            <p className="text-muted-foreground text-sm">{job.company}</p>
          </CardHeader>
          <CardContent className="grid gap-3">
            <p className="text-sm leading-relaxed whitespace-pre-line">
              {job.description}
            </p>
            <p className="text-muted-foreground text-xs">
              Posted by {job.postedBy.name}
            </p>
          </CardContent>
        </Card>
      )}

      {token ? (
        <ReferralForm jobId={jobId} />
      ) : (
        <p className="text-muted-foreground text-sm">
          Sign in from the dashboard to refer a candidate.
        </p>
      )}

      <section className="grid gap-3">
        <h2 className="text-lg font-semibold">
          Referrals{referrals ? ` (${referrals.length})` : ""}
        </h2>

        {!referrals && <Skeleton className="h-24 w-full" />}

        {referrals?.length === 0 && (
          <p className="text-muted-foreground text-sm">
            No referrals yet for this role.
          </p>
        )}

        {referrals?.map((referral) => (
          <ReferralCard key={referral.id} referral={referral} />
        ))}
      </section>
    </div>
  );
}
