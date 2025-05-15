import { createContext, useContext } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';

import { useAuth } from './auth-provider';
import type { FC, ReactNode } from 'react';
import type { Subscription } from '@/api/modules/subscriptions';
import { api } from '@/api';

type SubscriptionContextType = {
  subscription: Subscription | null;
  upgradeSubscription: (planId: string) => void;
  isUpgrading: boolean;
};

export const subscriptionContext = createContext<SubscriptionContextType>({
  subscription: null,
  upgradeSubscription: () => {},
  isUpgrading: false,
});

export const useSubscription = () => {
  const { subscription, upgradeSubscription, isUpgrading } =
    useContext(subscriptionContext);
  return {
    subscription,
    upgradeSubscription,
    isUpgrading,
  };
};

export const SubscriptionProvider: FC<{ children: ReactNode }> = (props) => {
  const { account } = useAuth();
  const { data: subscription, refetch } = useQuery({
    enabled: !!account,
    queryKey: ['subscription'],
    queryFn: () => api.subscriptions.getSubscription(),
  });

  const upgradeSubscription = useMutation({
    mutationFn: (planId: string) =>
      api.subscriptions.upgradeSubscription({
        data: { planId },
      }),
    onSuccess: () => {
      refetch();
    },
  });

  return (
    <subscriptionContext.Provider
      value={{
        subscription: subscription?.subscription ?? null,
        upgradeSubscription: upgradeSubscription.mutate,
        isUpgrading: upgradeSubscription.isPending,
      }}
    >
      {props.children}
    </subscriptionContext.Provider>
  );
};
