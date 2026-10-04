"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

import { useAuth } from "@/components/auth-provider";
import { useTopic } from "@/components/socket-provider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { ApiError, api } from "@/lib/api";
import type { Comment, WsEvent } from "@/lib/types";

function initials(name: string) {
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("");
}

function timeOf(iso: string) {
  return new Date(iso).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function CommentThread({ referralId }: { referralId: string }) {
  const { token, user } = useAuth();
  const [comments, setComments] = useState<Comment[] | null>(null);
  const [draft, setDraft] = useState("");
  const [busy, setBusy] = useState(false);
  const endRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let cancelled = false;
    api
      .listComments(referralId)
      .then((res) => {
        if (!cancelled) setComments(res.comments);
      })
      .catch(() => {
        if (!cancelled) setComments([]);
      });
    return () => {
      cancelled = true;
    };
  }, [referralId]);

  const onEvent = useCallback((event: WsEvent) => {
    if (event.type !== "comment.created") return;
    const incoming = event.payload;
    setComments((current) => {
      if (!current) return current;
      if (current.some((c) => c.id === incoming.id)) return current;
      return [...current, incoming];
    });
  }, []);

  useTopic(`referral:${referralId}`, onEvent);

  // Keep the newest message in view as the conversation grows.
  useEffect(() => {
    endRef.current?.scrollIntoView({ block: "nearest" });
  }, [comments?.length]);

  async function onSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (!token || !draft.trim()) return;

    setBusy(true);
    try {
      await api.createComment(token, referralId, draft);
      setDraft("");
    } catch (err) {
      toast.error(
        err instanceof ApiError ? err.message : "Could not post the comment",
      );
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="grid gap-3">
      <Separator />

      <div className="max-h-64 overflow-y-auto">
        {comments === null && (
          <p className="text-muted-foreground text-sm">Loading discussion</p>
        )}

        {comments?.length === 0 && (
          <p className="text-muted-foreground text-sm">
            No comments yet. Start the discussion.
          </p>
        )}

        <ul className="grid gap-3">
          {comments?.map((comment) => (
            <li key={comment.id} className="flex gap-2.5">
              <span className="bg-muted text-muted-foreground mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-full text-[11px] font-medium">
                {initials(comment.author.name)}
              </span>
              <div className="min-w-0">
                <p className="text-xs">
                  <span className="font-medium">
                    {comment.author.id === user?.id ? "You" : comment.author.name}
                  </span>
                  <span className="text-muted-foreground ml-2">
                    {timeOf(comment.createdAt)}
                  </span>
                </p>
                <p className="text-sm break-words">{comment.body}</p>
              </div>
            </li>
          ))}
        </ul>
        <div ref={endRef} />
      </div>

      {token ? (
        <form onSubmit={onSubmit} className="flex gap-2">
          <Input
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            placeholder="Add a comment"
            aria-label="Add a comment"
          />
          <Button type="submit" size="sm" disabled={busy || !draft.trim()}>
            Send
          </Button>
        </form>
      ) : (
        <p className="text-muted-foreground text-sm">Sign in to join the discussion.</p>
      )}
    </div>
  );
}
