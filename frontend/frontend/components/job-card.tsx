import Link from "next/link";
import { MessageSquare, Users } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Job } from "@/lib/types";

const statusVariant: Record<Job["status"], "default" | "secondary" | "outline"> = {
  open: "default",
  paused: "secondary",
  closed: "outline",
};

export function JobCard({ job }: { job: Job }) {
  return (
    <Card className="transition-colors hover:border-foreground/20">
      <CardHeader className="gap-1">
        <div className="flex flex-wrap items-start gap-2">
          <CardTitle className="text-base">
            <Link href={`/jobs/${job.id}`} className="hover:underline">
              {job.title}
            </Link>
          </CardTitle>
          <Badge variant={statusVariant[job.status]} className="ml-auto shrink-0">
            {job.status}
          </Badge>
        </div>
        <p className="text-muted-foreground text-sm">{job.company}</p>
      </CardHeader>

      <CardContent className="grid gap-3">
        <p className="text-sm leading-relaxed line-clamp-3">{job.description}</p>

        <div className="text-muted-foreground flex flex-wrap items-center gap-x-4 gap-y-1 text-xs">
          <span>Posted by {job.postedBy.name}</span>
          <span className="flex items-center gap-1">
            <Users className="size-3.5" aria-hidden />
            {job.referralCount} referral{job.referralCount === 1 ? "" : "s"}
          </span>
          <Link
            href={`/jobs/${job.id}`}
            className="ml-auto flex items-center gap-1 hover:underline"
          >
            <MessageSquare className="size-3.5" aria-hidden />
            Open and discuss
          </Link>
        </div>
      </CardContent>
    </Card>
  );
}
