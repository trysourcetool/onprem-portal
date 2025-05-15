import * as api from '@/api/instance';

type SubscriptionStatus = '' | 'trial' | 'active' | 'canceled' | 'past_due';

export type Subscription = {
  id: string;
  userId: string;
  planId: string;
  status: SubscriptionStatus;
  stripeCustomerId: string;
  stripeSubscriptionId: string;
  trialStart: string;
  trialEnd: string;
  createdAt: string;
  updatedAt: string;
};

export const getSubscription = async () => {
  const res = await api.get<{
    subscription: Subscription;
  }>({ path: '/subscriptions', auth: true });

  return res;
};

export const upgradeSubscription = async (params: {
  data: {
    planId: string;
  };
}) => {
  const res = await api.post<{
    code: number;
    message: string;
  }>({
    path: '/subscriptions/upgrade',
    data: params.data,
    auth: true,
  });

  return res;
};

export const cancelSubscription = async () => {
  const res = await api.post({
    path: '/subscriptions/cancel',
    auth: true,
  });

  return res;
};
