import { createFileRoute } from '@tanstack/react-router';
import { useEffect } from 'react';
import clsx from 'clsx';
import { useQuery } from '@tanstack/react-query';
import type { FC } from 'react';

import { PageHeader } from '@/components/common/page-header';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useAuth } from '@/components/provider/auth-provider';
import { checkAccountExpiredDays, checkTrialExpiredDays } from '@/lib/account';
import { PricingCard } from '@/components/common/pricing-card';
import { useSubscription } from '@/components/provider/subscription-provider';
import { api } from '@/api';

const VerticalLine: FC<{
  className?: string;
}> = ({ className }) => {
  return (
    <svg
      width="2"
      height="32"
      viewBox="0 0 2 32"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={clsx(
        'absolute left-1/2 hidden -translate-x-1/2 md:block',
        className,
      )}
    >
      <path
        d="M1 0V32"
        stroke="#E4E4E7"
        stroke-width="2"
        stroke-dasharray="2 2"
      />
    </svg>
  );
};

const HorizontalLine: FC<{
  className?: string;
}> = ({ className }) => {
  return (
    <svg
      width="160"
      height="2"
      viewBox="0 0 160 2"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={clsx(
        'absolute bottom-[calc(100%+32px)] hidden -translate-y-1/2 md:block',
        className,
      )}
    >
      <path
        d="M0 1H160"
        stroke="#E4E4E7"
        stroke-width="2"
        stroke-dasharray="2 2"
      />
    </svg>
  );
};

export default function Index() {
  const { account } = useAuth();
  const { data: plans } = useQuery({
    enabled: !!account,
    queryKey: ['plans'],
    queryFn: () => api.plan.getPlans(),
  });

  console.log({ plans });

  const {
    subscription,
    upgradeSubscription,
    cancelSubscription,
    isUpgrading,
    isCancelling,
  } = useSubscription();
  const { setBreadcrumbsState } = useBreadcrumbs();

  console.log({ subscription });

  useEffect(() => {
    setBreadcrumbsState?.([{ label: 'Pricing', to: '/pricing' }]);
  }, [setBreadcrumbsState]);

  return (
    <div>
      <PageHeader label="Pricing" />
      <div className="flex w-screen flex-col gap-4 px-4 py-6 md:w-auto md:gap-8 md:px-6">
        {account && subscription && (
          <div className="flex justify-center">
            <div className="bg-accent relative flex flex-col gap-2.5 rounded-md border p-10 text-center">
              <p className="font-bold">
                You're currently on a {checkTrialExpiredDays(subscription)}-day
                free trial
              </p>
              <p className="text-xs">
                You must select one of our plans within{' '}
                {checkAccountExpiredDays(account)} days.
              </p>
              <VerticalLine className="top-full" />
            </div>
          </div>
        )}

        <div className="flex flex-col items-stretch justify-center gap-6 px-20 py-8 md:flex-row">
          <div className="relative flex max-w-[296px] flex-1 flex-col [&>*]:flex-1">
            <VerticalLine className="bottom-full" />
            <HorizontalLine className="left-1/2" />
            <PricingCard
              title="Migrate to Community Edition"
              description="Perfect for testing and small projects"
              price={0}
              buttonLabel="Read documentation"
              onClick={() => {}}
              features={['Up to 5 users', 'Unlimited apps']}
              buttonDisabled={isUpgrading || isCancelling}
            />
          </div>
          <div className="relative flex max-w-[296px] flex-1 flex-col [&>*]:flex-1">
            <VerticalLine className="bottom-full" />
            <HorizontalLine className="left-1/2" />
            <HorizontalLine className="right-1/2" />
            <PricingCard
              title="Team"
              description="Ideal for growing teams and collaboration"
              price={10}
              period="user/month"
              buttonLabel="Upgrade"
              onClick={() => {}}
              features={['Staging Environment Available', 'More than 5 users']}
              isPopular
              buttonDisabled={isUpgrading || isCancelling}
            />
          </div>
          <div className="relative flex max-w-[296px] flex-1 flex-col [&>*]:flex-1">
            <VerticalLine className="bottom-full" />
            <HorizontalLine className="right-1/2" />
            <PricingCard
              title="Business"
              description="Advanced features for enterprise needs"
              price={24}
              period="month"
              buttonLabel="Upgrade"
              onClick={() => {}}
              features={[
                'Permission Control',
                'Unlimited Environments',
                'Audit Logs',
              ]}
              buttonDisabled={isUpgrading || isCancelling}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_default/pricing/')({
  component: Index,
});
