import { CommentThread } from "@/components/comment-thread";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Referral } from "@/lib/types";

export function ReferralCard({ referral }: { referral: Referral }) {
  return (
    <Card>
      <CardHeader className="gap-1">
        <div className="flex flex-wrap items-start gap-2">
          <CardTitle className="text-base">{referral.candidateName}</CardTitle>
          <Badge variant="secondary" className="ml-auto shrink-0">
            {referral.status}
          </Badge>
        </div>
        <p className="text-muted-foreground text-sm break-all">
          {referral.candidateEmail}
        </p>
      </CardHeader>

      <CardContent className="grid gap-3">
        {referral.message && (
          <p className="text-sm leading-relaxed">{referral.message}</p>
        )}

        <p className="text-muted-foreground text-xs">
          Referred by {referral.referrer.name}
        </p>

        <CommentThread referralId={referral.id} />
      </CardContent>
    </Card>
  );
}
