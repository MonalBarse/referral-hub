export type UserRef = {
  id: string;
  name: string;
  email: string;
};

export type User = {
  id: string;
  email: string;
  name: string;
  createdAt: string;
};

export type Job = {
  id: string;
  title: string;
  company: string;
  description: string;
  status: "open" | "paused" | "closed";
  postedBy: UserRef;
  referralCount: number;
  createdAt: string;
};

export type Referral = {
  id: string;
  jobId: string;
  candidateName: string;
  candidateEmail: string;
  message: string;
  status: string;
  referrer: UserRef;
  commentCount: number;
  createdAt: string;
};

export type Comment = {
  id: string;
  referralId: string;
  body: string;
  author: UserRef;
  createdAt: string;
};

export type WsEvent =
  | { type: "job.created"; topic: string; payload: Job; timestamp: string }
  | { type: "referral.created"; topic: string; payload: Referral; timestamp: string }
  | { type: "comment.created"; topic: string; payload: Comment; timestamp: string };
