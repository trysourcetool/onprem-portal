import { createFileRoute } from '@tanstack/react-router';
import { useCallback, useEffect } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';

import { toast } from 'sonner';
import { PageHeader } from '@/components/common/page-header';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useAuth } from '@/components/provider/auth-provider';
import { PricingCard } from '@/components/common/pricing-card';
import { useSubscription } from '@/components/provider/subscription-provider';
import { api } from '@/api';

export default function Index() {
  const { account } = useAuth();
  const { data: plans } = useQuery({
    enabled: !!account,
    queryKey: ['plans'],
    queryFn: () => api.plan.getPlans(),
  });

  const { subscription, upgradeSubscription, isUpgrading } = useSubscription();
  const { setBreadcrumbsState } = useBreadcrumbs();

  const createCheckoutSession = useMutation({
    mutationFn: (planId: string) =>
      api.stripe.createCheckoutSession({
        data: { planId },
      }),
    onSuccess: (data) => {
      window.location.href = data.url;
    },
    onError: (error) => {
      console.error(error);
      toast.error('Failed to create checkout session');
    },
  });

  const teamPlan = plans?.plans.find((plan) => plan.name === 'Team');
  const businessPlan = plans?.plans.find((plan) => plan.name === 'Business');

  useEffect(() => {
    setBreadcrumbsState?.([{ label: 'Pricing', to: '/settings/billing' }]);
  }, [setBreadcrumbsState]);

  const handleUpgrade = useCallback(
    (planId: string) => {
      if (subscription?.planId === planId) {
        return;
      }
      if (subscription?.plan) {
        upgradeSubscription(planId);
      } else {
        createCheckoutSession.mutate(planId);
      }
    },
    [
      upgradeSubscription,
      createCheckoutSession,
      subscription?.planId,
      subscription?.plan,
    ],
  );
  return (
    <div>
      <PageHeader label="Pricing" />
      <div className="flex w-screen flex-col gap-4 px-4 py-6 md:w-auto md:gap-8 md:px-6">
        <div className="flex flex-col items-stretch justify-center gap-6 px-20 py-8 md:flex-row">
          <div className="relative flex max-w-[296px] flex-1 flex-col [&>*]:flex-1">
            <PricingCard
              title="Migrate to Community Edition"
              description="Perfect for testing and small projects"
              price={0}
              buttonType="docs"
              features={['Up to 5 users', 'Unlimited apps']}
              buttonDisabled={isUpgrading}
            />
          </div>
          <div className="relative flex max-w-[296px] flex-1 flex-col [&>*]:flex-1">
            {teamPlan && (
              <PricingCard
                title={teamPlan.name}
                description="Ideal for growing teams and collaboration"
                price={teamPlan.price}
                period="user/month"
                buttonType="dialog"
                features={[
                  'Staging Environment Available',
                  'More than 5 users',
                ]}
                isPopular
                isCurrentPlan={subscription?.planId === teamPlan.id}
                buttonDisabled={isUpgrading}
                onDialogSubmit={() => handleUpgrade(teamPlan.id)}
              />
            )}
          </div>
          <div className="relative flex max-w-[296px] flex-1 flex-col [&>*]:flex-1">
            {businessPlan && (
              <PricingCard
                title={businessPlan.name}
                description="Advanced features for enterprise needs"
                price={businessPlan.price}
                period="month"
                buttonType="dialog"
                features={[
                  'Permission Control',
                  'Unlimited Environments',
                  'Audit Logs',
                ]}
                isCurrentPlan={subscription?.planId === businessPlan.id}
                buttonDisabled={isUpgrading}
                onDialogSubmit={() => handleUpgrade(businessPlan.id)}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_default/settings/billing/')({
  component: Index,
});
